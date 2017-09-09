// replication-manager - Replication Manager Monitoring and CLI for MariaDB and MySQL
// Authors: Guillaume Lefranc <guillaume@signal18.io>
//          Stephane Varoqui  <stephane@mariadb.com>
// This source code is licensed under the GNU General Public License, version 3.

package regtest

import (
	"time"

	"github.com/signal18/replication-manager/cluster"
)

func testSwitchoverAllSlavesDelayMultimasterNoRplChecksNoSemiSync(cluster *cluster.Cluster, conf string, test *cluster.Test) bool {
	cluster.SetMultiMaster(true)
	cluster.SetRplMaxDelay(8)
	cluster.SetRplChecks(false)
	err := cluster.DisableSemisync()
	if err != nil {
		cluster.LogPrintf("ERROR", "%s", err)
		return false
	}
	err = cluster.StopSlaves()
	if err != nil {
		cluster.LogPrintf("ERROR", "%s", err)
		return false
	}
	time.Sleep(15 * time.Second)
	SaveMasterURL := cluster.GetMaster().URL
	cluster.LogPrintf("TEST", "Master is %s", cluster.GetMaster().URL)
	cluster.SwitchoverWaitTest()
	cluster.LogPrintf("TEST", "New Master  %s ", cluster.GetMaster().URL)
	err = cluster.StartSlaves()
	if err != nil {
		cluster.LogPrintf("ERROR", "%s", err)
		return false
	}
	err = cluster.EnableSemisync()
	if err != nil {
		cluster.LogPrintf("ERROR", "%s", err)
		return false
	}
	time.Sleep(2 * time.Second)
	if cluster.GetMaster().URL == SaveMasterURL {
		cluster.LogPrintf("ERROR", "Saved Prefered master %s <>  from saved %s  ", SaveMasterURL, cluster.GetMaster().URL)
		return false
	}
	return true
}