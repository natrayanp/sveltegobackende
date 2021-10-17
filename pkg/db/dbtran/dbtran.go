package dbtran

import (
	"context"
	"fmt"
	"reflect"

	//_ "github.com/jackc/pgx/v4"
	//"github.com/jmoiron/sqlx"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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
	//Exec(query string, args ...interface{}) (sql.Result, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)

	//Preparex(query string) (*sqlx.Stmt, error)
	Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error)

	//Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)

	//QueryRowx(query string, args ...interface{}) *sqlx.Row
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row

	//Select(dest interface{}, query string, args ...interface{}) error
	//Select(ctx context.Context, db pgxscan.Querier, dst interface{}, query string, args ...interface{}) error
}

// A Txfn is a function that will be called with an initialized `Transaction` object
// that can be used for executing statements and queries against a database.
type TxFn func(context.Context, TranType, *pgxpool.Pool, Transaction) error

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn`
func WithTransaction(ctx context.Context, typ TranType, db *pgxpool.Pool, tra pgx.Tx, fn TxFn) (pgx.Tx, error) {
	var errc error

	if typ != TranTypeNoTran {

		if typ == TranTypeFullSet || typ == TranTypeFirstSet {
			tra, errc = db.Begin(ctx)
			if errc != nil {
				return nil, errc
			}
		} else if tra != nil {
			//txc = tra
		}

		defer func() {
			if p := recover(); p != nil {
				// a panic occurred, rollback and repanic
				tra.Rollback(ctx)
				panic(p)
			} else if errc != nil {
				// something went wrong, rollback
				tra.Rollback(ctx)
			} else {
				// all good, commit
				//But only if it is fullset or Lastset
				if typ == TranTypeFullSet || typ == TranTypeLastSet {
					errc = tra.Commit(ctx)
				}
			}
		}()
	}

	errc = fn(ctx, typ, db, tra)
	fmt.Println("+++++++++++++++++++++$$$end3")
	return tra, errc
}

// Result set of non select queries
type Resultset struct {
	RowsAffected int64
}

// A PipelineStmt is a simple wrapper for creating a statement consisting of
// a query and a set of arguments to be passed to that query.
type PipelineStmt struct {
	querytype    string
	query        string
	Resultstruct interface{}
	Resulterror  error
	args         []interface{}
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

func (ps *PipelineStmt) Exec(ctx context.Context, typ TranType, db *pgxpool.Pool, tx Transaction) error {
	var ct pgconn.CommandTag
	var err error

	if typ != TranTypeNoTran {
		ct, err = tx.Exec(ctx, ps.query, ps.args...)
	} else {
		ct, err = db.Exec(ctx, ps.query, ps.args...)
	}
	fmt.Println("Peieiei")

	dd := (ps.Resultstruct).(*Resultset)
	fmt.Println(ps.Resultstruct)
	if err != nil {

		dd.RowsAffected = -1
		return err
	}
	dd.RowsAffected = ct.RowsAffected()
	return nil
}

func (ps *PipelineStmt) Selects(ctx context.Context, typ TranType, db *pgxpool.Pool, tx Transaction) error {
	var rows pgx.Rows
	var err error
	fmt.Println("+++++ selects +++++")
	fmt.Println(typ)
	fmt.Println(reflect.TypeOf(ps.Resultstruct).Elem())
	if typ != TranTypeNoTran {
		fmt.Println("+++++++++++++++++++++qq1")
		rows, err = tx.Query(ctx, ps.query, ps.args...)

		fmt.Println("printrow:", rows)
		fmt.Println("error:", err)
		fmt.Println("+++++++++++++++++++++qq1")
	} else {
		fmt.Println("+++++++++++++++++++++qq")
		rows, err = db.Query(ctx, ps.query, ps.args...)
		fmt.Println("printrow:", rows)
		fmt.Println("error:", err)
		fmt.Println("+++++++++++++++++++++qq")
	}

	if err != nil {
		return err
	}
	fmt.Println("+++++++++++++++++++++$$$1")
	err = pgxscan.ScanAll(ps.Resultstruct, rows)
	fmt.Println(ps.Resultstruct)
	fmt.Println(err)
	//fmt.Println(len(ps.reultstruct))
	//fmt.Println(ps.reultstruct[0])
	//fmt.Println(ps.reultstruct[0]["companyid"])
	//fmt.Println(reflect.TypeOf(ps.reultstruct[0]["companyid"]))
	fmt.Println("+++++++++++++++++++++$$$2")
	if err != nil {
		return err
	}
	fmt.Println("+++++++++++++++++++++$$$end")

	return nil
}

// Runs the supplied statements within the transaction. If any statement fails, the transaction
// is rolled back, and the original error is returned.
//
// The `LastInsertId` from the previous statement will be passed to `Exec`. The zero-value (0) is
// used initially.
func RunPipeline(ctx context.Context, typ TranType, db *pgxpool.Pool, tx Transaction, stmts ...*PipelineStmt) error {
	var err error
	//var lastInsId int64
	var ps *PipelineStmt

	for _, ps = range stmts {
		fmt.Println("+++++++++++++++++++++$$$end1s")
		fmt.Println(ps.querytype)
		if ps.querytype != "select" {
			if ps.Resultstruct == nil {
				fmt.Println("assigned resulstset")
				ps.Resultstruct = &Resultset{}
			}
			err = ps.Exec(ctx, typ, db, tx)
			ps.Resulterror = err
		} else if ps.querytype == "select" {
			err = ps.Selects(ctx, typ, db, tx)
			fmt.Println(err)
			fmt.Println(ps)
			ps.Resulterror = err
			fmt.Println("+++++++++++++++++++++$$$end1")
		}

		if err != nil {
			return err
		}
	}
	fmt.Println("+++++++++++++++++++++$$$end2")
	return nil
}

/*
func RunSolo(db pgxpool.Pool, stmt *PipelineStmt) error {
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
