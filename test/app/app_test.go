// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build integration

package app

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/require"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/devkit"
	"openpitrix.io/openpitrix/pkg/devkit/opapp"
	"openpitrix.io/openpitrix/test/client/app_manager"
	"openpitrix.io/openpitrix/test/client/attachment_service"
	"openpitrix.io/openpitrix/test/models"
	"openpitrix.io/openpitrix/test/testutil"
)

var clientConfig = testutil.GetClientConfig()
var testTmpDir = testutil.GetTmpDir()

var Service = []string{"hyperpitrix", "openpitrix-am-service"}

const Vmbased = "vmbased"

func getTestIcon(t *testing.T) strfmt.Base64 {
	b, err := ioutil.ReadFile("testdata/logo.png")
	testutil.NoError(t, err, Service)
	return strfmt.Base64(b)
}

func getTestIcon1(t *testing.T) strfmt.Base64 {
	b, err := ioutil.ReadFile("testdata/logo1.png")
	testutil.NoError(t, err, Service)
	return strfmt.Base64(b)
}

func testAppIcon(t *testing.T, app *models.OpenpitrixApp) {
	iconAttachmentId := app.Icon
	filename := "raw"
	client := testutil.GetClient(clientConfig)

	getReq := attachment_service.NewGetAttachmentParams()
	getReq.SetAttachmentID(&iconAttachmentId)
	getReq.SetFilename(&filename)
	res, err := client.AttachmentService.GetAttachment(getReq, nil)
	testutil.NoError(t, err, Service)
	require.Equal(t, getTestIcon(t), res.Payload.Content)

	uploadAppAttachmentParams := app_manager.NewUploadAppAttachmentParams()
	uploadAppAttachmentParams.WithBody(
		&models.OpenpitrixUploadAppAttachmentRequest{
			AppID:             app.AppID,
			Type:              models.OpenpitrixUploadAppAttachmentRequestTypeIcon,
			AttachmentContent: getTestIcon1(t),
		})
	uploadAppAttachment, err := client.AppManager.UploadAppAttachment(uploadAppAttachmentParams, nil)
	testutil.NoError(t, err, Service)
	t.Log(uploadAppAttachment)

	getReq = attachment_service.NewGetAttachmentParams()
	getReq.SetAttachmentID(&iconAttachmentId)
	getReq.SetFilename(&filename)
	res, err = client.AttachmentService.GetAttachment(getReq, nil)
	testutil.NoError(t, err, Service)
	require.Equal(t, getTestIcon1(t), res.Payload.Content)
}

func preparePackage(t *testing.T, v string) strfmt.Base64 {
	var testAppName = "e2e_test_app"

	cfile := &opapp.Metadata{
		Name:        testAppName,
		Description: "An OpenPitrix app",
		Version:     v,
		AppVersion:  "1.0",
		ApiVersion:  devkit.ApiVersionV1,
	}

	os.MkdirAll(testTmpDir, 0755)
	_, err := devkit.Create(cfile, testTmpDir)

	testutil.NoError(t, err, Service)

	ch, err := devkit.LoadDir(path.Join(testTmpDir, testAppName))

	testutil.NoError(t, err, Service)

	name, err := devkit.Save(ch, testTmpDir)

	testutil.NoError(t, err, Service)

	t.Logf("save [%s] success", name)

	content, err := ioutil.ReadFile(name)

	testutil.NoError(t, err, Service)

	require.NoError(t, os.RemoveAll(testTmpDir))

	return strfmt.Base64(content)
}

func testVersionPackage(t *testing.T, appId string) {
	client := testutil.GetClient(clientConfig)

	modifyAppParams := app_manager.NewModifyAppParams()
	modifyAppParams.WithBody(
		&models.OpenpitrixModifyAppRequest{
			AppID: appId,
		})
	_, err := client.AppManager.ModifyApp(modifyAppParams, nil)

	testutil.NoError(t, err, Service)

	createAppVersionParams := app_manager.NewCreateAppVersionParams()
	createAppVersionParams.WithBody(
		&models.OpenpitrixCreateAppVersionRequest{
			AppID:   appId,
			Type:    Vmbased,
			Package: preparePackage(t, "0.0.2"),
		})
	createAppVersionResp, err := client.AppManager.CreateAppVersion(createAppVersionParams, nil)

	testutil.NoError(t, err, Service)

	versionId1 := createAppVersionResp.Payload.VersionID

	modifyAppVersionParams := app_manager.NewModifyAppVersionParams()
	modifyAppVersionParams.WithBody(
		&models.OpenpitrixModifyAppVersionRequest{
			VersionID: versionId1,
			Package:   preparePackage(t, "0.0.3"),
		})
	_, err = client.AppManager.ModifyAppVersion(modifyAppVersionParams, nil)

	testutil.NoError(t, err, Service)

	modifyAppVersionParams = app_manager.NewModifyAppVersionParams()
	modifyAppVersionParams.WithBody(
		&models.OpenpitrixModifyAppVersionRequest{
			VersionID: versionId1,
			Package:   preparePackage(t, "0.0.4"),
		})
	_, err = client.AppManager.ModifyAppVersion(modifyAppVersionParams, nil)

	testutil.NoError(t, err, Service)

	createAppVersionParams = app_manager.NewCreateAppVersionParams()
	createAppVersionParams.WithBody(
		&models.OpenpitrixCreateAppVersionRequest{
			AppID:   appId,
			Type:    Vmbased,
			Package: preparePackage(t, "0.1.0"),
		})
	createAppVersionResp, err = client.AppManager.CreateAppVersion(createAppVersionParams, nil)

	testutil.NoError(t, err, Service)

	versionId2 := createAppVersionResp.Payload.VersionID

	modifyAppVersionParams = app_manager.NewModifyAppVersionParams()
	modifyAppVersionParams.WithBody(
		&models.OpenpitrixModifyAppVersionRequest{
			VersionID: versionId2,
			Package:   preparePackage(t, "0.0.4"),
		})
	_, err = client.AppManager.ModifyAppVersion(modifyAppVersionParams, nil)

	testutil.NoError(t, err, Service)

	deleteAppVersionParams := app_manager.NewDeleteAppVersionParams()
	deleteAppVersionParams.WithBody(
		&models.OpenpitrixDeleteAppVersionRequest{
			VersionID: versionId2,
		})
	_, err = client.AppManager.DeleteAppVersion(deleteAppVersionParams, nil)

	testutil.NoError(t, err, Service)

	deleteAppVersionParams = app_manager.NewDeleteAppVersionParams()
	deleteAppVersionParams.WithBody(
		&models.OpenpitrixDeleteAppVersionRequest{
			VersionID: versionId1,
		})
	_, err = client.AppManager.DeleteAppVersion(deleteAppVersionParams, nil)

	testutil.NoError(t, err, Service)
}

func testVersionLifeCycle(t *testing.T, versionId string) {
	client := testutil.GetClient(clientConfig)

	modifyAppVersionParams := app_manager.NewModifyAppVersionParams()
	modifyAppVersionParams.WithBody(
		&models.OpenpitrixModifyAppVersionRequest{
			VersionID: versionId,
			Name:      "test_version2",
		})
	_, err := client.AppManager.ModifyAppVersion(modifyAppVersionParams, nil)

	testutil.NoError(t, err, Service)

	submitAppVersionParams := app_manager.NewSubmitAppVersionParams()
	submitAppVersionParams.WithBody(
		&models.OpenpitrixSubmitAppVersionRequest{
			VersionID: versionId,
		})
	_, err = client.AppManager.SubmitAppVersion(submitAppVersionParams, nil)

	testutil.NoError(t, err, Service)

	isvReviewAppVersionParams := app_manager.NewIsvReviewAppVersionParams()
	isvReviewAppVersionParams.WithBody(&models.OpenpitrixReviewAppVersionRequest{
		VersionID: versionId,
	})
	_, err = client.AppManager.IsvReviewAppVersion(isvReviewAppVersionParams, nil)
	testutil.NoError(t, err, Service)

	rejectAppVersionParams := app_manager.NewIsvRejectAppVersionParams()
	rejectAppVersionParams.WithBody(&models.OpenpitrixRejectAppVersionRequest{
		VersionID: versionId,
		Message:   "test message",
	})
	_, err = client.AppManager.IsvRejectAppVersion(rejectAppVersionParams, nil)

	testutil.NoError(t, err, Service)

	_, err = client.AppManager.SubmitAppVersion(submitAppVersionParams, nil)
	testutil.NoError(t, err, Service)

	isvReviewAppVersionParams = app_manager.NewIsvReviewAppVersionParams()
	isvReviewAppVersionParams.WithBody(&models.OpenpitrixReviewAppVersionRequest{
		VersionID: versionId,
	})
	_, err = client.AppManager.IsvReviewAppVersion(isvReviewAppVersionParams, nil)
	testutil.NoError(t, err, Service)

	passAppVersionParams := app_manager.NewIsvPassAppVersionParams()
	passAppVersionParams.WithBody(&models.OpenpitrixPassAppVersionRequest{
		VersionID: versionId,
	})
	_, err = client.AppManager.IsvPassAppVersion(passAppVersionParams, nil)
	testutil.NoError(t, err, Service)

	reviewAppVersionParams := app_manager.NewBusinessReviewAppVersionParams()
	reviewAppVersionParams.WithBody(&models.OpenpitrixReviewAppVersionRequest{
		VersionID: versionId,
	})
	_, err = client.AppManager.BusinessReviewAppVersion(reviewAppVersionParams, nil)
	testutil.NoError(t, err, Service)

	busPassAppVersionParams := app_manager.NewBusinessPassAppVersionParams()
	busPassAppVersionParams.WithBody(&models.OpenpitrixPassAppVersionRequest{
		VersionID: versionId,
	})
	_, err = client.AppManager.BusinessPassAppVersion(busPassAppVersionParams, nil)
	testutil.NoError(t, err, Service)

	devReviewAppVersionParams := app_manager.NewTechnicalReviewAppVersionParams()
	devReviewAppVersionParams.WithBody(&models.OpenpitrixReviewAppVersionRequest{
		VersionID: versionId,
	})
	_, err = client.AppManager.TechnicalReviewAppVersion(devReviewAppVersionParams, nil)
	testutil.NoError(t, err, Service)

	devPassAppVersionParams := app_manager.NewTechnicalPassAppVersionParams()
	devPassAppVersionParams.WithBody(&models.OpenpitrixPassAppVersionRequest{
		VersionID: versionId,
	})
	_, err = client.AppManager.TechnicalPassAppVersion(devPassAppVersionParams, nil)
	testutil.NoError(t, err, Service)

	releaseAppVersionParams := app_manager.NewReleaseAppVersionParams()
	releaseAppVersionParams.WithBody(
		&models.OpenpitrixReleaseAppVersionRequest{
			VersionID: versionId,
		})
	_, err = client.AppManager.ReleaseAppVersion(releaseAppVersionParams, nil)

	testutil.NoError(t, err, Service)

	suspendAppVersionParams := app_manager.NewSuspendAppVersionParams()
	suspendAppVersionParams.WithBody(
		&models.OpenpitrixSuspendAppVersionRequest{
			VersionID: versionId,
		})
	_, err = client.AppManager.SuspendAppVersion(suspendAppVersionParams, nil)

	testutil.NoError(t, err, Service)

	deleteAppVersionParams := app_manager.NewDeleteAppVersionParams()
	deleteAppVersionParams.WithBody(
		&models.OpenpitrixDeleteAppVersionRequest{
			VersionID: versionId,
		})
	_, err = client.AppManager.DeleteAppVersion(deleteAppVersionParams, nil)

	testutil.NoError(t, err, Service)
}

func testStatistics(t *testing.T) {
	client := testutil.GetClient(clientConfig)
	getStatisticsResp, err := client.AppManager.GetAppStatistics(nil, nil)
	testutil.NoError(t, err, Service)
	require.NotEmpty(t, getStatisticsResp.Payload.LastTwoWeekCreated)
	require.NotEmpty(t, getStatisticsResp.Payload.TopTenRepos)
	require.NotEmpty(t, getStatisticsResp.Payload.AppCount)
	require.NotEmpty(t, getStatisticsResp.Payload.RepoCount)
}

func TestApp(t *testing.T) {
	client := testutil.GetClient(clientConfig)

	// delete old app
	testAppName := "e2e_test_app"
	describeParams := app_manager.NewDescribeAppsParams()
	describeParams.SetName([]string{testAppName})
	describeParams.SetStatus([]string{constants.StatusDraft, constants.StatusActive})
	describeResp, err := client.AppManager.DescribeApps(describeParams, nil)

	testutil.NoError(t, err, Service)

	apps := describeResp.Payload.AppSet
	for _, app := range apps {
		deleteParams := app_manager.NewDeleteAppsParams()
		deleteParams.WithBody(
			&models.OpenpitrixDeleteAppsRequest{
				AppID: []string{app.AppID},
			})
		_, err := client.AppManager.DeleteApps(deleteParams, nil)

		testutil.NoError(t, err, Service)
	}
	// create app
	createParams := app_manager.NewCreateAppParams()
	createParams.WithBody(
		&models.OpenpitrixCreateAppRequest{
			Name:           testAppName,
			VersionPackage: preparePackage(t, "0.0.1"),
			VersionType:    Vmbased,
			Icon:           getTestIcon(t),
		})
	createResp, err := client.AppManager.CreateApp(createParams, nil)

	testutil.NoError(t, err, Service)

	appId := createResp.Payload.AppID
	versionId := createResp.Payload.VersionID
	// modify app

	//modifyParams := app_manager.NewModifyAppParams()
	//modifyParams.WithBody(
	//	&models.OpenpitrixModifyAppRequest{
	//		AppID:      appId,
	//		CategoryID: "aa,bb,cc,xx",
	//	})
	//modifyResp, err := client.AppManager.ModifyApp(modifyParams, nil)
	//
	//testutil.NoError(t, err, Service)
	//
	//t.Log(modifyResp)

	// describe app
	describeParams.WithAppID([]string{appId})
	describeResp, err = client.AppManager.DescribeApps(describeParams, nil)

	testutil.NoError(t, err, Service)

	apps = describeResp.Payload.AppSet

	require.Equal(t, 1, len(apps))

	testAppIcon(t, apps[0])
	testVersionPackage(t, appId)
	testVersionLifeCycle(t, versionId)
	testStatistics(t)

	// delete app
	deleteParams := app_manager.NewDeleteAppsParams()
	deleteParams.WithBody(&models.OpenpitrixDeleteAppsRequest{
		AppID: []string{appId},
	})
	deleteResp, err := client.AppManager.DeleteApps(deleteParams, nil)

	testutil.NoError(t, err, Service)

	t.Log(deleteResp)
	// describe deleted app
	describeParams.WithAppID([]string{appId})
	describeParams.WithStatus([]string{constants.StatusDeleted})
	describeParams.WithName(nil)
	describeResp, err = client.AppManager.DescribeApps(describeParams, nil)

	testutil.NoError(t, err, Service)

	apps = describeResp.Payload.AppSet

	require.Equal(t, 1, len(apps))

	app := apps[0]

	require.Equal(t, appId, app.AppID)

	require.Equal(t, constants.StatusDeleted, app.Status)

	t.Log("test app finish, all test is ok")
}
