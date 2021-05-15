package dbtran

import (
	"database/sql"

	_ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
)

//https://pseudomuto.com/2018/01/clean-sql-transactions-in-golang/

type TranType struct {
	t string
}

var (
	TranTypeFullSet  = TranType{"fullset"}
	TranTypeFirstSet = TranType{"firstset"}
	TranTypeMidSet   = TranType{"midset"}
	TranTypeLastSet  = TranType{"lastset"}
	TranTypeNoTran   = TranType{"notran"}
)

// Transaction is an interface that models the standard transaction in
// `database/sql`.
//
// To ensure `TxFn` funcs cannot commit or rollback a transaction (which is
// handled by `WithTransaction`), those methods are not included here.
type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Preparex(query string) (*sqlx.Stmt, error)
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Select(dest interface{}, query string, args ...interface{}) error
}

// A Txfn is a function that will be called with an initialized `Transaction` object
// that can be used for executing statements and queries against a database.
type TxFn func(TranType, *sqlx.DB, Transaction) error

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn`
func WithTransaction(typ TranType, db *sqlx.DB, tra *sqlx.Tx, fn TxFn) (*sqlx.Tx, error) {
	var txc *sqlx.Tx
	var errc error

	if typ != TranTypeNoTran {

		if typ == TranTypeFullSet || typ == TranTypeFirstSet {
			tx, err := db.Beginx()
			errc = err
			txc = tx
			if err != nil {
				return nil, err
			}
		} else if tra != nil {
			txc = tra
		}

		defer func() {
			if p := recover(); p != nil {
				// a panic occurred, rollback and repanic
				txc.Rollback()
				panic(p)
			} else if errc != nil {
				// something went wrong, rollback
				txc.Rollback()
			} else {
				// all good, commit
				//But only if it is fullset or Lastset
				if typ == TranTypeFullSet || typ == TranTypeLastSet {
					errc = txc.Commit()
				}
			}
		}()
	}

	errc = fn(typ, db, txc)
	return txc, errc
}

// A PipelineStmt is a simple wrapper for creating a statement consisting of
// a query and a set of arguments to be passed to that query.
type PipelineStmt struct {
	querytype   string
	query       string
	reultstruct interface{}
	resulterror error
	args        []interface{}
}

func NewPipelineStmt(querytype string, query string, reultstruct interface{}, args ...interface{}) *PipelineStmt {
	return &PipelineStmt{querytype, query, reultstruct, nil, args}
}

// Executes the statement within supplied transaction. The literal string `{LAST_INS_ID}`
// will be replaced with the supplied value to make chaining `PipelineStmt` objects together
// simple.
/*
func (ps *PipelineStmt) Exec(tx Transaction, lastInsertId int64) (sql.Result, error) {
	query := strings.Replace(ps.query, "{LAST_INS_ID}", strconv.Itoa(int(lastInsertId)), -1)
	return tx.Exec(query, ps.args...)
}
*/

func (ps *PipelineStmt) Exec(typ TranType, db *sqlx.DB, tx Transaction) (sql.Result, error) {
	if typ != TranTypeNoTran {
		return db.Exec(ps.query, ps.args...)
	} else {
		return tx.Exec(ps.query, ps.args...)
	}
}

func (ps *PipelineStmt) Selects(typ TranType, db *sqlx.DB, tx Transaction) error {
	if typ != TranTypeNoTran {
		return db.Select(ps.reultstruct, ps.query, ps.args...)
	} else {
		return tx.Select(ps.reultstruct, ps.query, ps.args...)
	}
}

// Runs the supplied statements within the transaction. If any statement fails, the transaction
// is rolled back, and the original error is returned.
//
// The `LastInsertId` from the previous statement will be passed to `Exec`. The zero-value (0) is
// used initially.
func RunPipeline(typ TranType, db *sqlx.DB, tx Transaction, stmts ...*PipelineStmt) (sql.Result, error) {
	var res sql.Result
	var err error
	//var lastInsId int64
	var ps *PipelineStmt

	for _, ps = range stmts {
		if ps.querytype != "select" {
			res, err = ps.Exec(typ, db, tx)
			ps.resulterror = err
		} else if ps.querytype == "select" {
			err = ps.Selects(typ, db, tx)
			ps.resulterror = err
		}

		if err != nil {
			return nil, err
		}
		/*
			if ps.querytype != "select" {
				lastInsId, err = res.LastInsertId()
				if err != nil {
					return nil, err
				}
			}
		*/
	}

	if ps.querytype != "select" {
		return res, nil
	} else {
		return nil, nil
	}
}

/*
func RunSolo(db sqlx.DB, stmt *PipelineStmt) error {
	ps := stmt

	err := ps.Selects(db)
	if err != nil {
		return err
	}

	return nil
}
*/

/*
Implentation example

func main() {
	db, err := sql.Open("VENDOR_HERE", "YOUR_DSN_HERE")
	handleError(err)

	defer db.Close()

		stmts := []*PipelineStmt{
		NewPipelineStmt("INSERT INTO table1(name) VALUES(?)", "some name"),
		NewPipelineStmt("INSERT INTO table2(table1_id, name) VALUES({LAST_INS_ID}, ?)", "other name"),
	}

	err = WithTransaction(db, func(typ TranType,db *sqlx.Tx,tx Transaction) error {
		_, err := RunPipeline(typ,db,tx, stmts...)
		return err
	})

	stmts := *PipelineStmt{
		NewPipelineStmt("SELECT * FROM AC.MYTABLE")
	}

		err = NoTran(db, func(db, tx Transaction) error {
		_, err := RunSolo(db, stmts)
		return err
	})

	handleError(err)
	log.Println("Done.")
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

*/
