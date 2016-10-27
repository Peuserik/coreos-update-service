package coreos

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/Sirupsen/logrus"
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb/util"
	"time"
	"strconv"
)

type Storage interface {
	RegisterNode(machineId string, track string, oem string, version string)
	ListNodes() map[string]string
	UpdateVersion(coreOSVersion CoreOSVersion)
	GetVersion(version string) CoreOSVersion
	LogEvent(machineId string, etype string, eResult string, eError string)
	UpdateTracks(map[string] string)
	GetTracks() map[string] string
}

// Provides a level db implementation for the record storage
type LocalDB struct {
	Db *leveldb.DB
}

const nodes_space = 0x00
const events_space = 0x01
const versions_space  = 0x02
const tracks_space = 0x03

func NewLocalDB(file string) Storage {
	db, err := leveldb.OpenFile(file, nil)
	if (err != nil){
		logrus.Panic(err)
	}
	return LocalDB{Db: db}
}


func buildKey(space byte, key []byte) []byte{
	result := []byte {space}
	for i := 0; i < len(key); i++{
		result = append(result, key[i])
	}
	return result
}


func (ldb LocalDB) RegisterNode(machineId string, track string, oem string, version string) {
	value, _ := json.Marshal(map[string] string {
		"machineId": machineId,
		"track": track,
		"oem": oem,
		"version": version,
	})
	if err := ldb.Db.Put(buildKey(nodes_space,[]byte(machineId)), value, nil); err != nil {
		logrus.Error(err)
	}
}

func (ldb LocalDB) ListNodes() map[string]string{
	hosts := map[string]string {}
	iter := ldb.Db.NewIterator(&util.Range{Start: []byte{nodes_space}, Limit: []byte{events_space}},nil)
	for iter.Next() {
		hosts[string(iter.Key())] = string(iter.Value())
	}
	iter.Release()
	return hosts
}

func (ldb LocalDB) UpdateVersion(coreOsVersion CoreOSVersion){
	key := buildKey(versions_space, []byte(coreOsVersion.VersionId))
	data, _:= json.Marshal(coreOsVersion)
	logrus.Infof("Insert version %s", coreOsVersion.VersionId)
	if err := ldb.Db.Put(key, data, nil); err != nil {
		logrus.Error(err)
	}
}

func (ldb LocalDB) GetVersion(version string) CoreOSVersion {
	key := buildKey(versions_space, []byte(version))
	if value, err := ldb.Db.Get(key, nil); err != nil {
		logrus.Error(err)
	} else {
		var coreosVersion CoreOSVersion
		json.Unmarshal(value, &coreosVersion)
		return coreosVersion
	}
	// To do return the error
	return CoreOSVersion{}
}

func (ldb LocalDB) LogEvent(machineId string, etype string, eResult string, eError string)  {
	key := buildKey(events_space, []byte(strconv.FormatInt(time.Now().UnixNano(),10)))

	data := map [string] string {
		"TimeStamp": time.Now().Format(time.RFC3339),
		"MachineId": machineId,
		"Type": etype,
		"Result": eResult,
		"Error": eError,
	}
	dataBytes, _ := json.Marshal(data)
	if err :=  ldb.Db.Put(key, dataBytes, nil); err != nil {
		logrus.Error(err)
	}
}

func (ldb LocalDB) UpdateTracks(tracks map[string] string) {
	key := buildKey(tracks_space, []byte("tracks"))
	data, _ := json.Marshal(tracks)
	if err :=  ldb.Db.Put(key, data, nil); err != nil {
		logrus.Error(err)
	}
}


func (ldb LocalDB) GetTracks() map[string] string {
	key := buildKey(tracks_space, []byte("tracks"))

	if value, err := ldb.Db.Get(key, nil); err != nil {
		logrus.Error(err)
	}else {
		var tracks map [string]string
		json.Unmarshal(value,&tracks)
		return tracks
	}
	return nil
}


