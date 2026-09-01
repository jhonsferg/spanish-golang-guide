// Ejemplo del Capítulo 64: SQL crudo de bajo nivel — configuración del
// connection pool, prepared statements reutilizados, y transacciones
// con rollback automático ante error.
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Configuración del pool: crítico en producción contra una base real
	// (aquí con SQLite en memoria el efecto es mínimo, pero la API es la
	// misma sin importar el driver).
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if _, err := db.Exec(`CREATE TABLE cuentas (
		id      INTEGER PRIMARY KEY,
		titular TEXT NOT NULL,
		saldo   INTEGER NOT NULL
	)`); err != nil {
		log.Fatal(err)
	}

	// Prepared statement: se compila una vez y se reutiliza en cada
	// llamada — más rápido que Exec() repetido con el mismo SQL.
	insertar, err := db.Prepare("INSERT INTO cuentas (id, titular, saldo) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer insertar.Close()

	for _, c := range []struct {
		id      int
		titular string
		saldo   int
	}{
		{1, "Ana", 500},
		{2, "Beto", 200},
	} {
		if _, err := insertar.Exec(c.id, c.titular, c.saldo); err != nil {
			log.Fatal(err)
		}
	}

	// Transferencia como transacción: o se aplican AMBOS cambios, o
	// ninguno — nunca un estado intermedio inconsistente.
	if err := transferir(context.Background(), db, 1, 2, 150); err != nil {
		log.Fatal(err)
	}
	imprimirSaldos(db)

	// Transferencia que debe fallar (fondos insuficientes) y hacer rollback.
	err = transferir(context.Background(), db, 2, 1, 999999)
	if err != nil {
		fmt.Println("transferencia rechazada (esperado):", err)
	}
	imprimirSaldos(db) // los saldos deben seguir iguales: el rollback funcionó
}

var ErrFondosInsuficientes = errors.New("fondos insuficientes")

func transferir(ctx context.Context, db *sql.DB, origenID, destinoID, monto int) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // no-op si ya hubo Commit

	var saldoOrigen int
	if err := tx.QueryRowContext(ctx, "SELECT saldo FROM cuentas WHERE id = ?", origenID).Scan(&saldoOrigen); err != nil {
		return err
	}
	if saldoOrigen < monto {
		return ErrFondosInsuficientes
	}

	if _, err := tx.ExecContext(ctx, "UPDATE cuentas SET saldo = saldo - ? WHERE id = ?", monto, origenID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "UPDATE cuentas SET saldo = saldo + ? WHERE id = ?", monto, destinoID); err != nil {
		return err
	}

	return tx.Commit()
}

func imprimirSaldos(db *sql.DB) {
	rows, err := db.Query("SELECT id, titular, saldo FROM cuentas ORDER BY id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, saldo int
		var titular string
		rows.Scan(&id, &titular, &saldo)
		fmt.Printf("  cuenta %d (%s): %d\n", id, titular, saldo)
	}
}
