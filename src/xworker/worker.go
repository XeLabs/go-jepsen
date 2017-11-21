/*
 * benchyou
 * xelabs.org
 *
 * Copyright (c) XeLabs
 * GPL License
 *
 */

package xworker

import (
	"fmt"
	"log"
	"xcommon"

	"github.com/XeLabs/go-mysqlstack/driver"
)

type Metric struct {
	WNums  uint64
	WCosts uint64
	WMax   uint64
	WMin   uint64
	QNums  uint64
	QCosts uint64
	QMax   uint64
	QMin   uint64
	QErrs  uint64
}

type Worker struct {
	// session
	S driver.Conn

	// mertric
	M *Metric

	// engine
	E string

	// xid
	XID string
}

func CreateWorkers(conf *xcommon.Conf, threads int) []Worker {
	var workers []Worker
	for i := 0; i < threads; i++ {
		conn, err := driver.NewConn(conf.Mysql_user, conf.Mysql_password, fmt.Sprintf("%s:%d", conf.Mysql_host, conf.Mysql_port), "", "utf8")
		if err != nil {
			log.Panicf("create.worker.error:%v", err)
		}
		workers = append(workers, Worker{
			S: conn,
			M: &Metric{},
			E: conf.Mysql_table_engine,
		},
		)
	}
	return workers
}

func AllWorkersMetric(workers []Worker) *Metric {
	all := &Metric{}
	for _, worker := range workers {
		all.WNums += worker.M.WNums
		all.WCosts += worker.M.WCosts
		all.QErrs += worker.M.QErrs

		if all.WMax < worker.M.WMax {
			all.WMax = worker.M.WMax
		}

		if all.WMin > worker.M.WMin {
			all.WMin = worker.M.WMin
		}

		all.QNums += worker.M.QNums
		all.QCosts += worker.M.QCosts

		if all.QMax < worker.M.QMax {
			all.QMax = worker.M.QMax
		}

		if all.QMin > worker.M.QMin {
			all.QMin = worker.M.QMin
		}
	}

	return all
}

func StopWorkers(workers []Worker) {
	for _, worker := range workers {
		worker.S.Close()
	}
}
