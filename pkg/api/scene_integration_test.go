//go:build integration
// +build integration

package api_test

import (
	"reflect"
	"testing"

	"github.com/stashapp/stash-box/pkg/models"
	"github.com/stashapp/stash-box/pkg/user"
)

type sceneTestRunner struct {
	testRunner
}

func createSceneTestRunner(t *testing.T) *sceneTestRunner {
	return &sceneTestRunner{
		testRunner: *asModify(t),
	}
}

func (s *sceneTestRunner) testCreateScene() {
	title := "Title"
	details := "Details"
	date := "2003-02-01"

	performer, _ := s.createTestPerformer(nil)
	studio, _ := s.createTestStudio(nil)
	tag, _ := s.createTestTag(nil)

	performerID := performer.ID
	studioID := studio.ID
	tagID := tag.ID

	performerAlias := "alias"

	input := models.SceneCreateInput{
		Title:   &title,
		Details: &details,
		Date:    &date,
		Fingerprints: []*models.FingerprintEditInput{
			s.generateSceneFingerprint(nil),
		},
		StudioID: &studioID,
		Performers: []*models.PerformerAppearanceInput{
			{
				PerformerID: performerID,
				As:          &performerAlias,
			},
		},
		Urls: []*models.URLInput{
			{
				URL:  "URL",
				Type: "Type",
			},
		},
		TagIds: []string{
			tagID,
		},
	}

	scene, err := s.client.createScene(input)

	if err != nil {
		s.t.Errorf("Error creating scene: %s", err.Error())
		return
	}

	s.verifyCreatedScene(input, scene)
}

func comparePerformers(input []*models.PerformerAppearanceInput, performers []*performerAppearance) bool {
	if len(performers) != len(input) {
		return false
	}

	for i, v := range performers {
		performerID := v.Performer.ID
		if performerID != input[i].PerformerID {
			return false
		}

		if v.As != input[i].As {
			if v.As == nil || input[i].As == nil {
				return false
			}

			if *v.As != *input[i].As {
				return false
			}
		}
	}

	return true
}

func comparePerformersInput(input, performers []*models.PerformerAppearanceInput) bool {
	if len(performers) != len(input) {
		return false
	}

	for i, v := range performers {
		performerID := v.PerformerID
		if performerID != input[i].PerformerID {
			return false
		}

		if v.As != input[i].As {
			if v.As == nil || input[i].As == nil {
				return false
			}

			if *v.As != *input[i].As {
				return false
			}
		}
	}

	return true
}

func compareTags(tagIDs []string, tags []*idObject) bool {
	if len(tags) != len(tagIDs) {
		return false
	}

	for i, v := range tags {
		tagID := v.ID
		if tagID != tagIDs[i] {
			return false
		}
	}

	return true
}

func compareFingerprints(input []*models.FingerprintEditInput, fingerprints []*fingerprint) bool {
	if len(input) != len(fingerprints) {
		return false
	}

	for i, v := range fingerprints {
		if input[i].Algorithm != v.Algorithm || input[i].Hash != v.Hash {
			return false
		}
	}

	return true
}

func compareFingerprintsInput(input, fingerprints []*models.FingerprintEditInput) bool {
	if len(input) != len(fingerprints) {
		return false
	}

	for i, v := range fingerprints {
		if input[i].Algorithm != v.Algorithm || input[i].Hash != v.Hash {
			return false
		}
	}

	return true
}

func (s *sceneTestRunner) verifyCreatedScene(input models.SceneCreateInput, scene *sceneOutput) {
	// ensure basic attributes are set correctly
	if scene.ID == "" {
		s.t.Errorf("Expected created scene id to be non-zero")
	}

	if !reflect.DeepEqual(scene.Title, input.Title) {
		s.fieldMismatch(*input.Title, scene.Title, "Title")
	}

	if !reflect.DeepEqual(scene.Details, input.Details) {
		s.fieldMismatch(input.Details, scene.Details, "Details")
	}

	// ensure urls were set correctly
	if !compareUrls(input.Urls, scene.Urls) {
		s.fieldMismatch(input.Urls, scene.Urls, "Urls")
	}

	if !reflect.DeepEqual(scene.Date, input.Date) {
		s.fieldMismatch(*input.Date, scene.Date, "Date")
	}

	if !compareFingerprints(input.Fingerprints, scene.Fingerprints) {
		s.fieldMismatch(input.Fingerprints, scene.Fingerprints, "Fingerprints")
	}

	if !comparePerformers(input.Performers, scene.Performers) {
		s.fieldMismatch(input.Performers, scene.Performers, "Performers")
	}

	if !compareTags(input.TagIds, scene.Tags) {
		s.fieldMismatch(input.TagIds, scene.Tags, "Tags")
	}
}

func (s *sceneTestRunner) testFindSceneById() {
	createdScene, err := s.createTestScene(nil)
	if err != nil {
		return
	}

	scene, err := s.client.findScene(createdScene.ID)

	if err != nil {
		s.t.Errorf("Error finding scene: %s", err.Error())
		return
	}

	// ensure returned scene is not nil
	if scene == nil {
		s.t.Error("Did not find scene by id")
		return
	}

	// ensure values were set
	if *createdScene.Title != *scene.Title {
		s.fieldMismatch(createdScene.Title, scene.Title, "Title")
	}
}

func (s *sceneTestRunner) testFindSceneByFingerprint() {
	createdScene, err := s.createTestScene(nil)
	if err != nil {
		s.t.Errorf("Error creating scene: %s", err.Error())
		return
	}

	fingerprints := createdScene.Fingerprints
	if err != nil {
		s.t.Errorf("Error finding scene: %s", err.Error())
		return
	}
	fingerprint := models.FingerprintQueryInput{
		Algorithm: fingerprints[0].Algorithm,
		Hash:      fingerprints[0].Hash,
	}

	scenes, err := s.client.findSceneByFingerprint(fingerprint)
	if err != nil {
		s.t.Errorf("Error finding scene: %s", err.Error())
		return
	}

	// ensure returned scene is not nil
	if len(scenes) == 0 {
		s.t.Error("Did not find scene by fingerprint")
		return
	}

	// ensure values were set
	if *createdScene.Title != *scenes[0].Title {
		s.fieldMismatch(createdScene.Title, scenes[0].Title, "Title")
	}
}

func (s *sceneTestRunner) testFindScenesByFingerprints() {
	scene1Title := "asdasd"
	scene1Input := models.SceneCreateInput{
		Title: &scene1Title,
		Fingerprints: []*models.FingerprintEditInput{
			s.generateSceneFingerprint(nil),
		},
	}
	createdScene1, err := s.createTestScene(&scene1Input)
	if err != nil {
		return
	}
	createdScene2, err := s.createTestScene(nil)
	if err != nil {
		return
	}

	fingerprintList := []string{}
	fingerprints := createdScene1.Fingerprints
	fingerprintList = append(fingerprintList, fingerprints[0].Hash)
	fingerprints = createdScene2.Fingerprints
	fingerprintList = append(fingerprintList, fingerprints[0].Hash)

	scenes, err := s.client.findScenesByFingerprints(fingerprintList)
	if err != nil {
		s.t.Errorf("Error finding scenes: %s", err.Error())
		return
	}

	// ensure only two scenes are returned
	if len(scenes) != 2 {
		s.t.Error("Did not get correct amount of scenes by fingerprint")
		return
	}

	// ensure values were set
	if *createdScene1.Title != *scenes[0].Title {
		s.fieldMismatch(createdScene1.Title, scenes[0].Title, "Title")
	}
	if *createdScene2.Title != *scenes[1].Title {
		s.fieldMismatch(createdScene2.Title, scenes[1].Title, "Title")
	}
}

func (s *sceneTestRunner) testUpdateScene() {
	title := "Title"
	details := "Details"
	date := "2003-02-01"

	performer, _ := s.createTestPerformer(nil)
	studio, _ := s.createTestStudio(nil)
	tag, _ := s.createTestTag(nil)

	performerID := performer.ID
	studioID := studio.ID
	tagID := tag.ID

	performerAlias := "alias"

	input := models.SceneCreateInput{
		Title:   &title,
		Details: &details,
		Date:    &date,
		Fingerprints: []*models.FingerprintEditInput{
			// fingerprint that will be kept
			s.generateSceneFingerprint([]string{
				userDB.none.ID.String(),
				userDB.admin.ID.String(),
			}),
			// fingerprint that will be removed
			s.generateSceneFingerprint(nil),
		},
		StudioID: &studioID,
		Performers: []*models.PerformerAppearanceInput{
			{
				PerformerID: performerID,
				As:          &performerAlias,
			},
		},
		Urls: []*models.URLInput{
			{
				URL:  "URL",
				Type: "Type",
			},
		},
		TagIds: []string{
			tagID,
		},
	}

	createdScene, err := s.createTestScene(&input)
	if err != nil {
		return
	}

	sceneID := createdScene.ID

	newTitle := "NewTitle"
	newDetails := "NewDetails"
	newDate := "2001-02-03"

	performer, _ = s.createTestPerformer(nil)
	studio, _ = s.createTestStudio(nil)
	tag, _ = s.createTestTag(nil)

	performerID = performer.ID
	studioID = studio.ID
	tagID = tag.ID

	performerAlias = "updatedAlias"

	updateInput := models.SceneUpdateInput{
		ID:      sceneID,
		Title:   &newTitle,
		Details: &newDetails,
		Date:    &newDate,
		Fingerprints: []*models.FingerprintEditInput{
			input.Fingerprints[0],
			s.generateSceneFingerprint(nil),
		},
		Performers: []*models.PerformerAppearanceInput{
			{
				PerformerID: performerID,
				As:          &performerAlias,
			},
		},
		Urls: []*models.URLInput{
			{
				URL:  "URL",
				Type: "Type",
			},
		},
		StudioID: &studioID,
		TagIds: []string{
			tagID,
		},
	}

	scene, err := s.client.updateScene(updateInput)
	if err != nil {
		s.t.Errorf("Error updating scene: %s", err.Error())
		return
	}

	s.verifyUpdatedScene(updateInput, scene)

	// ensure fingerprint changes were enacted
	s.verifyUpdatedFingerprints(input.Fingerprints, updateInput.Fingerprints, scene)

	// ensure submissions count was maintained
	originalFP := input.Fingerprints[0]
	foundFP := false
	for _, f := range scene.Fingerprints {
		if originalFP.Algorithm == f.Algorithm && originalFP.Hash == f.Hash {
			foundFP = true
			if f.Submissions != 2 {
				s.t.Errorf("Incorrect fingerprint submissions count: %d", f.Submissions)
			}
		}
	}

	if !foundFP {
		s.t.Error("Could not find original fingerprint")
	}
}

func (s *sceneTestRunner) verifyUpdatedScene(input models.SceneUpdateInput, scene *sceneOutput) {
	// ensure basic attributes are set correctly
	if !reflect.DeepEqual(scene.Title, input.Title) {
		s.fieldMismatch(input.Title, scene.Title, "Title")
	}

	if !reflect.DeepEqual(scene.Details, input.Details) {
		s.fieldMismatch(input.Details, scene.Details, "Details")
	}

	if !reflect.DeepEqual(scene.Date, input.Date) {
		s.fieldMismatch(input.Date, scene.Date, "Date")
	}

	// ensure urls were set correctly
	if !compareUrls(input.Urls, scene.Urls) {
		s.fieldMismatch(input.Urls, scene.Urls, "Urls")
	}

	if !comparePerformers(input.Performers, scene.Performers) {
		s.fieldMismatch(input.Performers, scene.Performers, "Performers")
	}

	if !compareTags(input.TagIds, scene.Tags) {
		s.fieldMismatch(input.TagIds, scene.Tags, "Tags")
	}
}

func (s *sceneTestRunner) verifyUpdatedFingerprints(original, updated []*models.FingerprintEditInput, scene *sceneOutput) {
	hashExists := func(h *models.FingerprintEditInput, vs []*models.FingerprintEditInput) bool {
		for _, v := range vs {
			if h.Algorithm == v.Algorithm && h.Hash == v.Hash {
				return true
			}
		}

		return false
	}

	inOutput := func(h *models.FingerprintEditInput) bool {
		for _, hh := range scene.Fingerprints {
			if hh.Algorithm == h.Algorithm && hh.Hash == h.Hash {
				return true
			}
		}

		return false
	}

	for _, o := range original {
		// find in updated
		if hashExists(o, updated) {
			// exists, so ensure hash exists in output
			if !inOutput(o) {
				s.t.Errorf("existing hash %s missing in output", o.Hash)
			}
		} else {
			// not exists, ensure not in output
			if inOutput(o) {
				s.t.Errorf("removed hash %s still in output", o.Hash)
			}
		}
	}

	for _, u := range updated {
		// find in original
		if !hashExists(u, original) {
			// new hash, ensure in output
			if !inOutput(u) {
				s.t.Errorf("new hash %s missing in output", u.Hash)
			}
		}
	}
}

func (s *sceneTestRunner) testDestroyScene() {
	createdScene, err := s.createTestScene(nil)
	if err != nil {
		return
	}

	sceneID := createdScene.ID

	destroyed, err := s.client.destroyScene(models.SceneDestroyInput{
		ID: sceneID,
	})

	if err != nil {
		s.t.Errorf("Error destroying scene: %s", err.Error())
		return
	}

	if !destroyed {
		s.t.Error("Scene was not destroyed")
		return
	}

	// ensure cannot find scene
	foundScene, err := s.client.findScene(sceneID)
	if err != nil {
		s.t.Errorf("Error finding scene after destroying: %s", err.Error())
		return
	}

	if foundScene != nil {
		s.t.Error("Found scene after destruction")
	}

	// TODO - ensure scene was not removed
}

func (s *sceneTestRunner) testSubmitFingerprint() {
	createdScene, err := s.createTestScene(nil)
	if err != nil {
		return
	}

	fp := s.generateSceneFingerprint(nil)

	if _, err := s.client.submitFingerprint(models.FingerprintSubmission{
		SceneID: createdScene.ID,
		Fingerprint: &models.FingerprintInput{
			Hash:      fp.Hash,
			Algorithm: fp.Algorithm,
			Duration:  fp.Duration,
		},
	}); err != nil {
		s.t.Errorf("Error submitting fingerprint: %s", err.Error())
		return
	}

	scene, err := s.client.findScene(createdScene.ID)
	if err != nil {
		s.t.Errorf("Error finding scene: %s", err.Error())
		return
	}

	// verify created fingerprint
	expected := fingerprint{
		Hash:        fp.Hash,
		Algorithm:   fp.Algorithm,
		Duration:    fp.Duration,
		Submissions: 1,
	}
	actualFP := scene.Fingerprints[1]
	actual := fingerprint{
		Hash:        actualFP.Hash,
		Algorithm:   actualFP.Algorithm,
		Duration:    actualFP.Duration,
		Submissions: actualFP.Submissions,
	}
	if !reflect.DeepEqual(actual, expected) {
		s.fieldMismatch(expected, *scene.Fingerprints[1], "fingerprints")
	}

	// submit the same fingerprint - should not add and should not error
	if _, err := s.client.submitFingerprint(models.FingerprintSubmission{
		SceneID: createdScene.ID,
		Fingerprint: &models.FingerprintInput{
			Hash:      fp.Hash,
			Algorithm: fp.Algorithm,
			Duration:  fp.Duration,
		},
	}); err != nil {
		s.t.Errorf("Error submitting fingerprint: %s", err.Error())
		return
	}
}

func (s *sceneTestRunner) testSubmitFingerprintUnmatch() {
	createdScene, err := s.createTestScene(nil)
	if err != nil {
		return
	}

	unmatch := true
	if _, err := s.client.submitFingerprint(models.FingerprintSubmission{
		SceneID: createdScene.ID,
		Fingerprint: &models.FingerprintInput{
			Hash:      createdScene.Fingerprints[0].Hash,
			Algorithm: createdScene.Fingerprints[0].Algorithm,
			Duration:  createdScene.Fingerprints[0].Duration,
		},
		Unmatch: &unmatch,
	}); err != nil {
		s.t.Errorf("Error submitting fingerprint: %s", err.Error())
		return
	}

	scene, err := s.client.findScene(createdScene.ID)
	if err != nil {
		s.t.Errorf("Error finding scene: %s", err.Error())
		return
	}

	if len(scene.Fingerprints) > 0 {
		s.fieldMismatch([]*fingerprint{}, scene.Fingerprints, "fingerprints")
	}
}

func (s *sceneTestRunner) testSubmitFingerprintModify() {
	createdScene, err := s.createTestScene(nil)
	if err != nil {
		return
	}

	fp := s.generateSceneFingerprint(nil)

	if _, err := s.client.submitFingerprint(models.FingerprintSubmission{
		SceneID: createdScene.ID,
		Fingerprint: &models.FingerprintInput{
			Hash:      fp.Hash,
			Algorithm: fp.Algorithm,
			Duration:  fp.Duration,
			UserIds: []string{
				userDB.edit.ID.String(),
				userDB.none.ID.String(),
				userDB.read.ID.String(),
			},
		},
	}); err != nil {
		s.t.Errorf("Error submitting fingerprint: %s", err.Error())
		return
	}

	scene, err := s.client.findScene(createdScene.ID)
	if err != nil {
		s.t.Errorf("Error finding scene: %s", err.Error())
		return
	}

	// verify created fingerprint
	expected := fingerprint{
		Hash:        fp.Hash,
		Algorithm:   fp.Algorithm,
		Duration:    fp.Duration,
		Submissions: 3,
	}
	actualFP := scene.Fingerprints[0]
	actual := fingerprint{
		Hash:        actualFP.Hash,
		Algorithm:   actualFP.Algorithm,
		Duration:    actualFP.Duration,
		Submissions: actualFP.Submissions,
	}
	if !reflect.DeepEqual(actual, expected) {
		s.fieldMismatch(expected, *scene.Fingerprints[0], "fingerprints")
	}

	// submit the same fingerprint - should add
	if _, err := s.client.submitFingerprint(models.FingerprintSubmission{
		SceneID: createdScene.ID,
		Fingerprint: &models.FingerprintInput{
			Hash:      fp.Hash,
			Algorithm: fp.Algorithm,
			Duration:  fp.Duration,
		},
	}); err != nil {
		s.t.Errorf("Error submitting fingerprint: %s", err.Error())
		return
	}

	scene, err = s.client.findScene(createdScene.ID)
	if err != nil {
		s.t.Errorf("Error finding scene: %s", err.Error())
		return
	}

	expected.Submissions = 4
	actual.Submissions = scene.Fingerprints[0].Submissions

	if !reflect.DeepEqual(actual, expected) {
		s.fieldMismatch(expected, *scene.Fingerprints[0], "fingerprints")
	}
}

func (s *sceneTestRunner) testSubmitFingerprintUnmatchModify() {
	createdScene, err := s.createTestScene(nil)
	if err != nil {
		return
	}

	fp := s.generateSceneFingerprint(nil)

	if _, err := s.client.submitFingerprint(models.FingerprintSubmission{
		SceneID: createdScene.ID,
		Fingerprint: &models.FingerprintInput{
			Hash:      fp.Hash,
			Algorithm: fp.Algorithm,
			Duration:  fp.Duration,
			UserIds: []string{
				userDB.edit.ID.String(),
				userDB.none.ID.String(),
				userDB.read.ID.String(),
			},
		},
	}); err != nil {
		s.t.Errorf("Error submitting fingerprint: %s", err.Error())
		return
	}

	unmatch := true
	if _, err := s.client.submitFingerprint(models.FingerprintSubmission{
		SceneID: createdScene.ID,
		Fingerprint: &models.FingerprintInput{
			Hash:      fp.Hash,
			Algorithm: fp.Algorithm,
			Duration:  fp.Duration,
			UserIds: []string{
				userDB.edit.ID.String(),
			},
		},
		Unmatch: &unmatch,
	}); err != nil {
		s.t.Errorf("Error submitting fingerprint: %s", err.Error())
		return
	}

	scene, err := s.client.findScene(createdScene.ID)
	if err != nil {
		s.t.Errorf("Error finding scene: %s", err.Error())
		return
	}

	expected := fingerprint{
		Hash:        fp.Hash,
		Algorithm:   fp.Algorithm,
		Duration:    fp.Duration,
		Submissions: 2,
	}
	actualFP := scene.Fingerprints[0]
	actual := fingerprint{
		Hash:        actualFP.Hash,
		Algorithm:   actualFP.Algorithm,
		Duration:    actualFP.Duration,
		Submissions: actualFP.Submissions,
	}
	if !reflect.DeepEqual(actual, expected) {
		s.fieldMismatch(expected, *scene.Fingerprints[0], "fingerprints")
	}
}

func (s *sceneTestRunner) verifyQueryScenesResult(filter models.SceneFilterType, ids []string) {
	s.t.Helper()

	page := 1
	pageSize := 10
	querySpec := models.QuerySpec{
		Page:    &page,
		PerPage: &pageSize,
	}

	results, err := s.client.queryScenes(&filter, &querySpec)
	if err != nil {
		s.t.Errorf("Error querying scenes: %s", err.Error())
		return
	}

	if results.Count != len(ids) {
		s.t.Errorf("Expected %d query result, got %d", len(ids), results.Count)
		return
	}

	for _, id := range ids {
		found := false
		for _, scene := range results.Scenes {
			if scene.ID == id {
				found = true
				break
			}
		}

		if !found {
			s.t.Errorf("Missing scene with ID %s, got %v", id, results.Scenes)
			return
		}
	}
}

func (s *sceneTestRunner) verifyInvalidModifier(filter models.SceneFilterType) {
	s.t.Helper()

	page := 1
	pageSize := 10
	querySpec := models.QuerySpec{
		Page:    &page,
		PerPage: &pageSize,
	}

	defer func() {
		if r := recover(); r != nil {
			// success
		} else {
			s.t.Error("Expected error for invalid modifier")
		}
	}()
	s.resolver.Query().QueryScenes(s.ctx, &filter, &querySpec)
}

func (s *sceneTestRunner) testQueryScenesByStudio() {
	studio1, _ := s.createTestStudio(nil)
	studio2, _ := s.createTestStudio(nil)

	studio1ID := studio1.ID
	studio2ID := studio2.ID

	prefix := "testQueryScenesByStudio_"
	scene1Title := prefix + "scene1Title"
	scene2Title := prefix + "scene2Title"
	scene3Title := prefix + "scene3Title"

	input := models.SceneCreateInput{
		StudioID: &studio1ID,
		Title:    &scene1Title,
	}

	scene1, err := s.createTestScene(&input)
	if err != nil {
		return
	}

	input.StudioID = &studio2ID
	input.Title = &scene2Title
	scene2, err := s.createTestScene(&input)
	if err != nil {
		return
	}

	input.StudioID = nil
	input.Title = &scene3Title
	scene3, err := s.createTestScene(&input)
	if err != nil {
		return
	}

	scene1ID := scene1.ID
	scene2ID := scene2.ID
	scene3ID := scene3.ID

	// test equals
	filter := models.SceneFilterType{
		Studios: &models.MultiIDCriterionInput{
			Value:    []string{studio1ID},
			Modifier: models.CriterionModifierEquals,
		},
	}

	s.verifyQueryScenesResult(filter, []string{scene1ID})

	filter.Studios.Modifier = models.CriterionModifierNotEquals
	filter.Title = &scene2Title
	s.verifyQueryScenesResult(filter, []string{scene2ID})

	filter.Studios.Modifier = models.CriterionModifierIsNull
	filter.Title = &scene3Title
	s.verifyQueryScenesResult(filter, []string{scene3ID})

	filter.Studios.Modifier = models.CriterionModifierNotNull
	filter.Title = &scene1Title
	s.verifyQueryScenesResult(filter, []string{scene1ID})

	filter.Studios.Modifier = models.CriterionModifierIncludes
	filter.Studios.Value = []string{studio1ID, studio2ID}
	filter.Title = nil
	s.verifyQueryScenesResult(filter, []string{scene1ID, scene2ID})

	filter.Studios.Modifier = models.CriterionModifierExcludes
	filter.Studios.Value = []string{studio1ID}
	filter.Title = &scene2Title
	s.verifyQueryScenesResult(filter, []string{scene2ID})

	// test invalid modifiers
	filter.Studios.Modifier = models.CriterionModifierGreaterThan
	s.verifyInvalidModifier(filter)

	filter.Studios.Modifier = models.CriterionModifierLessThan
	s.verifyInvalidModifier(filter)

	filter.Studios.Modifier = models.CriterionModifierIncludesAll
	s.verifyInvalidModifier(filter)
}

func (s *sceneTestRunner) testQueryScenesByPerformer() {
	performer1, _ := s.createTestPerformer(nil)
	performer2, _ := s.createTestPerformer(nil)

	performer1ID := performer1.ID
	performer2ID := performer2.ID

	prefix := "testQueryScenesByPerformer_"
	scene1Title := prefix + "scene1Title"
	scene2Title := prefix + "scene2Title"
	scene3Title := prefix + "scene3Title"

	input := models.SceneCreateInput{
		Performers: []*models.PerformerAppearanceInput{
			{
				PerformerID: performer1ID,
			},
		},
		Title: &scene1Title,
	}

	scene1, err := s.createTestScene(&input)
	if err != nil {
		return
	}

	input.Performers[0].PerformerID = performer2ID
	input.Title = &scene2Title
	scene2, err := s.createTestScene(&input)
	if err != nil {
		return
	}

	input.Performers = append(input.Performers, &models.PerformerAppearanceInput{
		PerformerID: performer1ID,
	})
	input.Title = &scene3Title
	scene3, err := s.createTestScene(&input)
	if err != nil {
		return
	}

	scene1ID := scene1.ID
	scene2ID := scene2.ID
	scene3ID := scene3.ID

	titleSearch := prefix
	filter := models.SceneFilterType{
		Performers: &models.MultiIDCriterionInput{
			Value:    []string{performer1ID},
			Modifier: models.CriterionModifierIncludes,
		},
		Title: &titleSearch,
	}

	s.verifyQueryScenesResult(filter, []string{scene1ID, scene3ID})

	filter.Performers.Modifier = models.CriterionModifierExcludes
	s.verifyQueryScenesResult(filter, []string{scene2ID})

	filter.Performers.Modifier = models.CriterionModifierIncludesAll
	filter.Performers.Value = append(filter.Performers.Value, performer2ID)
	s.verifyQueryScenesResult(filter, []string{scene3ID})

	// test invalid modifiers
	filter.Performers.Modifier = models.CriterionModifierGreaterThan
	s.verifyInvalidModifier(filter)

	filter.Performers.Modifier = models.CriterionModifierLessThan
	s.verifyInvalidModifier(filter)

	filter.Performers.Modifier = models.CriterionModifierEquals
	s.verifyInvalidModifier(filter)

	filter.Performers.Modifier = models.CriterionModifierNotEquals
	s.verifyInvalidModifier(filter)

	filter.Performers.Modifier = models.CriterionModifierIsNull
	s.verifyInvalidModifier(filter)

	filter.Performers.Modifier = models.CriterionModifierNotNull
	s.verifyInvalidModifier(filter)
}

func (s *sceneTestRunner) testQueryScenesByTag() {
	tag1, _ := s.createTestTag(nil)
	tag2, _ := s.createTestTag(nil)

	tag1ID := tag1.ID
	tag2ID := tag2.ID

	prefix := "testQueryScenesByTag_"
	scene1Title := prefix + "scene1Title"
	scene2Title := prefix + "scene2Title"
	scene3Title := prefix + "scene3Title"

	input := models.SceneCreateInput{
		TagIds: []string{
			tag1ID,
		},
		Title: &scene1Title,
	}

	scene1, err := s.createTestScene(&input)
	if err != nil {
		return
	}

	input.TagIds[0] = tag2ID
	input.Title = &scene2Title
	scene2, err := s.createTestScene(&input)
	if err != nil {
		return
	}

	input.TagIds = append(input.TagIds, tag1ID)
	input.Title = &scene3Title
	scene3, err := s.createTestScene(&input)
	if err != nil {
		return
	}

	scene1ID := scene1.ID
	scene2ID := scene2.ID
	scene3ID := scene3.ID

	titleSearch := prefix
	filter := models.SceneFilterType{
		Tags: &models.MultiIDCriterionInput{
			Value:    []string{tag1ID},
			Modifier: models.CriterionModifierIncludes,
		},
		Title: &titleSearch,
	}

	s.verifyQueryScenesResult(filter, []string{scene1ID, scene3ID})

	filter.Tags.Modifier = models.CriterionModifierExcludes
	s.verifyQueryScenesResult(filter, []string{scene2ID})

	filter.Tags.Modifier = models.CriterionModifierIncludesAll
	filter.Tags.Value = append(filter.Tags.Value, tag2ID)
	s.verifyQueryScenesResult(filter, []string{scene3ID})

	// test invalid modifiers
	filter.Tags.Modifier = models.CriterionModifierGreaterThan
	s.verifyInvalidModifier(filter)

	filter.Tags.Modifier = models.CriterionModifierLessThan
	s.verifyInvalidModifier(filter)

	filter.Tags.Modifier = models.CriterionModifierEquals
	s.verifyInvalidModifier(filter)

	filter.Tags.Modifier = models.CriterionModifierNotEquals
	s.verifyInvalidModifier(filter)

	filter.Tags.Modifier = models.CriterionModifierIsNull
	s.verifyInvalidModifier(filter)

	filter.Tags.Modifier = models.CriterionModifierNotNull
	s.verifyInvalidModifier(filter)
}

func (s *sceneTestRunner) testUnauthorisedSceneModify() {
	// test each api interface - all require modify so all should fail
	_, err := s.resolver.Mutation().SceneCreate(s.ctx, models.SceneCreateInput{})
	if err != user.ErrUnauthorized {
		s.t.Errorf("SceneCreate: got %v want %v", err, user.ErrUnauthorized)
	}

	_, err = s.resolver.Mutation().SceneUpdate(s.ctx, models.SceneUpdateInput{})
	if err != user.ErrUnauthorized {
		s.t.Errorf("SceneUpdate: got %v want %v", err, user.ErrUnauthorized)
	}

	_, err = s.resolver.Mutation().SceneDestroy(s.ctx, models.SceneDestroyInput{})
	if err != user.ErrUnauthorized {
		s.t.Errorf("SceneDestroy: got %v want %v", err, user.ErrUnauthorized)
	}
}

func (s *sceneTestRunner) testUnauthorisedSceneQuery() {
	// test each api interface - all require read so all should fail
	_, err := s.resolver.Query().FindScene(s.ctx, "")
	if err != user.ErrUnauthorized {
		s.t.Errorf("FindScene: got %v want %v", err, user.ErrUnauthorized)
	}

	_, err = s.resolver.Query().QueryScenes(s.ctx, nil, nil)
	if err != user.ErrUnauthorized {
		s.t.Errorf("QueryScenes: got %v want %v", err, user.ErrUnauthorized)
	}
}

func TestCreateScene(t *testing.T) {
	pt := createSceneTestRunner(t)
	pt.testCreateScene()
}

func TestFindSceneById(t *testing.T) {
	pt := createSceneTestRunner(t)
	pt.testFindSceneById()
}

func TestFindSceneByFingerprint(t *testing.T) {
	pt := createSceneTestRunner(t)
	pt.testFindSceneByFingerprint()
}

func TestFindScenesByFingerprints(t *testing.T) {
	pt := createSceneTestRunner(t)
	pt.testFindScenesByFingerprints()
}

func TestUpdateScene(t *testing.T) {
	pt := createSceneTestRunner(t)
	pt.testUpdateScene()
}

// TestUpdateSceneTitle is removed due to no longer allowing
// partial updates

func TestDestroyScene(t *testing.T) {
	pt := createSceneTestRunner(t)
	pt.testDestroyScene()
}

func TestQueryScenesByStudio(t *testing.T) {
	pt := createSceneTestRunner(t)
	pt.testQueryScenesByStudio()
}

func TestQueryScenesByPerformer(t *testing.T) {
	pt := createSceneTestRunner(t)
	pt.testQueryScenesByPerformer()
}

func TestQueryScenesByTag(t *testing.T) {
	pt := createSceneTestRunner(t)
	pt.testQueryScenesByTag()
}

func TestUnauthorisedSceneModify(t *testing.T) {
	pt := &sceneTestRunner{
		testRunner: *asRead(t),
	}
	pt.testUnauthorisedSceneModify()
}

func TestUnauthorisedSceneQuery(t *testing.T) {
	pt := &sceneTestRunner{
		testRunner: *asNone(t),
	}
	pt.testUnauthorisedSceneQuery()
}

func TestSubmitFingerprint(t *testing.T) {
	pt := createSceneTestRunner(t)
	pt.testSubmitFingerprint()
}

func TestSubmitFingerprintUnmatch(t *testing.T) {
	pt := createSceneTestRunner(t)
	pt.testSubmitFingerprintUnmatch()
}

func TestSubmitFingerprintModify(t *testing.T) {
	pt := createSceneTestRunner(t)
	pt.testSubmitFingerprintModify()
}

func TestSubmitFingerprintUnmatchModify(t *testing.T) {
	pt := createSceneTestRunner(t)
	pt.testSubmitFingerprintUnmatchModify()
}
