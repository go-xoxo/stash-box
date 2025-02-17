//go:build integration
// +build integration

package api_test

import (
	"reflect"
	"testing"

	"github.com/stashapp/stash-box/pkg/models"
	"github.com/stashapp/stash-box/pkg/user"
)

type tagTestRunner struct {
	testRunner
}

func createTagTestRunner(t *testing.T) *tagTestRunner {
	return &tagTestRunner{
		testRunner: *asModify(t),
	}
}

func (s *tagTestRunner) testCreateTag() {
	description := "Description"

	input := models.TagCreateInput{
		Name:        s.generateTagName(),
		Description: &description,
	}

	tag, err := s.resolver.Mutation().TagCreate(s.ctx, input)

	if err != nil {
		s.t.Errorf("Error creating tag: %s", err.Error())
		return
	}

	s.verifyCreatedTag(input, tag)
}

func (s *tagTestRunner) verifyCreatedTag(input models.TagCreateInput, tag *models.Tag) {
	// ensure basic attributes are set correctly
	if input.Name != tag.Name {
		s.fieldMismatch(input.Name, tag.Name, "Name")
	}

	r := s.resolver.Tag()

	id, _ := r.ID(s.ctx, tag)
	if id == "" {
		s.t.Errorf("Expected created tag id to be non-zero")
	}

	if v, _ := r.Description(s.ctx, tag); !reflect.DeepEqual(v, input.Description) {
		s.fieldMismatch(*input.Description, v, "Description")
	}
}

func (s *tagTestRunner) testFindTagById() {
	createdTag, err := s.createTestTag(nil)
	if err != nil {
		return
	}

	tagID := createdTag.ID
	tag, err := s.resolver.Query().FindTag(s.ctx, &tagID, nil)
	if err != nil {
		s.t.Errorf("Error finding tag: %s", err.Error())
		return
	}

	// ensure returned tag is not nil
	if tag == nil {
		s.t.Error("Did not find tag by id")
		return
	}

	// ensure values were set
	if createdTag.Name != tag.Name {
		s.fieldMismatch(createdTag.Name, tag.Name, "Name")
	}
}

func (s *tagTestRunner) testFindTagByName() {
	createdTag, err := s.createTestTag(nil)
	if err != nil {
		return
	}

	tagName := createdTag.Name

	tag, err := s.resolver.Query().FindTag(s.ctx, nil, &tagName)
	if err != nil {
		s.t.Errorf("Error finding tag: %s", err.Error())
		return
	}

	// ensure returned tag is not nil
	if tag == nil {
		s.t.Error("Did not find tag by name")
		return
	}

	// ensure values were set
	if createdTag.Name != tag.Name {
		s.fieldMismatch(createdTag.Name, tag.Name, "Name")
	}
}

func (s *tagTestRunner) testUpdateTag() {
	createdTag, err := s.createTestTag(nil)
	if err != nil {
		return
	}

	tagID := createdTag.ID

	newDescription := "newDescription"

	updateInput := models.TagUpdateInput{
		ID:          tagID,
		Description: &newDescription,
	}

	updatedTag, err := s.resolver.Mutation().TagUpdate(s.ctx, updateInput)
	if err != nil {
		s.t.Errorf("Error updating tag: %s", err.Error())
		return
	}

	updateInput.Name = &createdTag.Name
	s.verifyUpdatedTag(updateInput, updatedTag)
}

func (s *tagTestRunner) verifyUpdatedTag(input models.TagUpdateInput, tag *models.Tag) {
	// ensure basic attributes are set correctly
	if input.Name != nil && *input.Name != tag.Name {
		s.fieldMismatch(*input.Name, tag.Name, "Name")
	}

	r := s.resolver.Tag()

	if v, _ := r.Description(s.ctx, tag); !reflect.DeepEqual(v, input.Description) {
		s.fieldMismatch(input.Description, v, "Description")
	}
}

func (s *tagTestRunner) testDestroyTag() {
	createdTag, err := s.createTestTag(nil)
	if err != nil {
		return
	}

	tagID := createdTag.ID

	destroyed, err := s.resolver.Mutation().TagDestroy(s.ctx, models.TagDestroyInput{
		ID: tagID,
	})
	if err != nil {
		s.t.Errorf("Error destroying tag: %s", err.Error())
		return
	}

	if !destroyed {
		s.t.Error("Tag was not destroyed")
		return
	}

	// ensure cannot find tag
	foundTag, err := s.resolver.Query().FindTag(s.ctx, &tagID, nil)
	if err != nil {
		s.t.Errorf("Error finding tag after destroying: %s", err.Error())
		return
	}

	if foundTag != nil {
		s.t.Error("Found tag after destruction")
	}

	// TODO - ensure scene was not removed
}

func (s *tagTestRunner) testUnauthorisedTagModify() {
	// test each api interface - all require modify so all should fail
	_, err := s.resolver.Mutation().TagCreate(s.ctx, models.TagCreateInput{})
	if err != user.ErrUnauthorized {
		s.t.Errorf("TagCreate: got %v want %v", err, user.ErrUnauthorized)
	}

	_, err = s.resolver.Mutation().TagUpdate(s.ctx, models.TagUpdateInput{})
	if err != user.ErrUnauthorized {
		s.t.Errorf("TagUpdate: got %v want %v", err, user.ErrUnauthorized)
	}

	_, err = s.resolver.Mutation().TagDestroy(s.ctx, models.TagDestroyInput{})
	if err != user.ErrUnauthorized {
		s.t.Errorf("TagDestroy: got %v want %v", err, user.ErrUnauthorized)
	}
}

func (s *tagTestRunner) testUnauthorisedTagQuery() {
	// test each api interface - all require read so all should fail
	_, err := s.resolver.Query().FindTag(s.ctx, nil, nil)
	if err != user.ErrUnauthorized {
		s.t.Errorf("FindTag: got %v want %v", err, user.ErrUnauthorized)
	}

	_, err = s.resolver.Query().QueryTags(s.ctx, nil, nil)
	if err != user.ErrUnauthorized {
		s.t.Errorf("QueryTags: got %v want %v", err, user.ErrUnauthorized)
	}
}

func TestCreateTag(t *testing.T) {
	pt := createTagTestRunner(t)
	pt.testCreateTag()
}

func TestFindTagById(t *testing.T) {
	pt := createTagTestRunner(t)
	pt.testFindTagById()
}

func TestFindTagByName(t *testing.T) {
	pt := createTagTestRunner(t)
	pt.testFindTagByName()
}

func TestUpdateTag(t *testing.T) {
	pt := createTagTestRunner(t)
	pt.testUpdateTag()
}

func TestDestroyTag(t *testing.T) {
	pt := createTagTestRunner(t)
	pt.testDestroyTag()
}

func TestUnauthorisedTagModify(t *testing.T) {
	pt := &tagTestRunner{
		testRunner: *asRead(t),
	}
	pt.testUnauthorisedTagModify()
}

func TestUnauthorisedTagQuery(t *testing.T) {
	pt := &tagTestRunner{
		testRunner: *asNone(t),
	}
	pt.testUnauthorisedTagQuery()
}
