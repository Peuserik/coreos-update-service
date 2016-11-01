package coreos

import (
	"github.com/Sirupsen/logrus"
	"github.com/coreos/go-omaha/omaha"
	"strconv"
)

const CoreOs = "{e96281a6-d1af-4bde-9a0a-97b76e56dc57}"

type CoreOsUpdater struct {
	Storage Storage
}

type CoreOSVersion struct {
	VersionId string
	URL       string
	Name      string
	Hash      string
	Signature string
	Size      int
}

func (cou CoreOsUpdater) noUpdateResponse() *omaha.Response {
	var response = omaha.NewResponse("update-controller")
	response.DayStart.ElapsedSeconds = "0"
	app := response.AddApp(CoreOs)
	app.Status = "ok"
	updateCheck := app.AddUpdateCheck()
	updateCheck.Status = "noupdate"
	return response
}

func (cou CoreOsUpdater) updateResponse(appRequest *omaha.App) *omaha.Response {
	tracks := cou.Storage.GetTracks()

	if version, ok := tracks[appRequest.Track]; ok {
		logrus.Infof("Requested track %s with target version %s", appRequest.Track, version)
		if appRequest.Version == version {
			logrus.Info("Version up to day")
			return cou.noUpdateResponse()

		} else {

			var response = omaha.NewResponse("update-controller")

			coreosVersion := cou.Storage.GetVersion(version)
			response.DayStart.ElapsedSeconds = "0"
			app := response.AddApp(CoreOs)
			app.Status = "ok"
			updateCheck := app.AddUpdateCheck()
			updateCheck.Status = "ok"
			updateCheck.AddUrl(coreosVersion.URL)

			manifest := updateCheck.AddManifest(coreosVersion.VersionId)
			manifest.AddPackage(coreosVersion.Hash, coreosVersion.Name, strconv.Itoa(coreosVersion.Size), true)
			action := manifest.AddAction("postinstall")
			action.Sha256 = coreosVersion.Signature
			action.NeedsAdmin = false
			action.IsDelta = false
			return response

		}
	} else {
		logrus.Errorf("Requested track not configured %s", appRequest.Track)
		return cou.noUpdateResponse()
	}
}

func (cou CoreOsUpdater) GetNodes() map[string]string {
	return cou.Storage.ListNodes()
}

func (cou CoreOsUpdater) UpdateVersion(coreOSVersion CoreOSVersion) {
	cou.Storage.UpdateVersion(coreOSVersion)
}

func (cou CoreOsUpdater) registerNode(app *omaha.App) {
	cou.Storage.RegisterNode(app.MachineID, app.Track, app.OEM, app.Version)
}

func (cou CoreOsUpdater) UpdateTracks(tracks map[string]string) {
	cou.Storage.UpdateTracks(tracks)
}

func (cou CoreOsUpdater) GetTracks() map[string]string {
	return cou.Storage.GetTracks()
}

func (cou CoreOsUpdater) CoreOsUpdate(request omaha.Request) *omaha.Response {
	app := request.Apps[0]

	cou.registerNode(app)

	// Log events
	for _, event := range app.Events {
		cou.Storage.LogEvent(app.MachineID, event.Type, event.Result, event.ErrorCode)
	}

	// Check for updates
	if app.UpdateCheck != nil {
		return cou.updateResponse(app)
	} else {
		return cou.noUpdateResponse()
	}
}
