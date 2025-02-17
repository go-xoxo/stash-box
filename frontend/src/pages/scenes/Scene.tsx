import React, { useContext } from "react";
import { Link, useHistory } from "react-router-dom";
import { Button, Card, Tabs, Tab, Table } from "react-bootstrap";
import { faCheckCircle } from "@fortawesome/free-solid-svg-icons";

import {
  Scene_findScene as Scene,
  Scene_findScene_fingerprints as Fingerprint,
} from "src/graphql/definitions/Scene";
import { useEdits, TargetTypeEnum, VoteStatusEnum } from "src/graphql";
import AuthContext from "src/AuthContext";
import {
  canEdit,
  getImage,
  getUrlByType,
  tagHref,
  performerHref,
  studioHref,
  createHref,
  formatDuration,
  formatDateTime,
  formatPendingEdits,
} from "src/utils";
import {
  ROUTE_SCENE_EDIT,
  ROUTE_SCENES,
  ROUTE_SCENE_DELETE,
} from "src/constants/route";
import {
  GenderIcon,
  TagLink,
  PerformerName,
  Icon,
} from "src/components/fragments";
import { EditList } from "src/components/list";

const DEFAULT_TAB = "description";

interface Props {
  scene: Scene;
}

const SceneComponent: React.FC<Props> = ({ scene }) => {
  const history = useHistory();
  const activeTab = history.location.hash?.slice(1) || DEFAULT_TAB;
  const auth = useContext(AuthContext);

  const { data: editData } = useEdits({
    filter: {
      per_page: 1,
    },
    editFilter: {
      target_type: TargetTypeEnum.SCENE,
      target_id: scene.id,
      status: VoteStatusEnum.PENDING,
    },
  });
  const pendingEditCount = editData?.queryEdits.count;

  const setTab = (tab: string | null) =>
    history.push({ hash: tab === DEFAULT_TAB ? "" : `#${tab}` });

  const performers = scene.performers
    .map((performance) => {
      const { performer } = performance;
      return (
        <Link
          key={performer.id}
          to={performerHref(performer)}
          className="scene-performer"
        >
          <GenderIcon gender={performer.gender} />
          <PerformerName performer={performer} as={performance.as} />
        </Link>
      );
    })
    .map((p, index) => (index % 2 === 2 ? [" • ", p] : p));

  function maybeRenderSubmitted(fingerprint: Fingerprint) {
    if (fingerprint.user_submitted) {
      return (
        <span className="user-submitted" title="Submitted by you">
          <Icon icon={faCheckCircle} />
        </span>
      );
    }
  }

  const fingerprints = scene.fingerprints.map((fingerprint) => (
    <tr key={fingerprint.hash}>
      <td>{fingerprint.algorithm}</td>
      <td>
        <Link
          to={`${createHref(ROUTE_SCENES)}?fingerprint=${fingerprint.hash}`}
        >
          {fingerprint.hash}
        </Link>
      </td>
      <td>
        <span title={`${fingerprint.duration}s`}>
          {formatDuration(fingerprint.duration)}
        </span>
      </td>
      <td>
        {fingerprint.submissions}
        {maybeRenderSubmitted(fingerprint)}
      </td>
      <td>{formatDateTime(fingerprint.created)}</td>
      <td>{formatDateTime(fingerprint.updated)}</td>
    </tr>
  ));
  const tags = [...scene.tags]
    .sort((a, b) => {
      if (a.name > b.name) return 1;
      if (a.name < b.name) return -1;
      return 0;
    })
    .map((tag) => (
      <li key={tag.name}>
        <TagLink title={tag.name} link={tagHref(tag)} />
      </li>
    ));

  return (
    <>
      <Card className="scene-info">
        <Card.Header>
          <div className="float-right">
            {canEdit(auth.user) && !scene.deleted && (
              <>
                <Link to={createHref(ROUTE_SCENE_EDIT, { id: scene.id })}>
                  <Button>Edit</Button>
                </Link>
                <Link
                  to={createHref(ROUTE_SCENE_DELETE, { id: scene.id })}
                  className="ml-2"
                >
                  <Button variant="danger">Delete</Button>
                </Link>
              </>
            )}
          </div>
          <h3>
            {scene.deleted ? (
              <del>{scene.title}</del>
            ) : (
              <span>{scene.title}</span>
            )}
          </h3>
          <h6>
            {scene.studio && (
              <>
                <Link to={studioHref(scene.studio)}>{scene.studio.name}</Link>
                <span className="mx-1">•</span>
              </>
            )}
            {scene.date}
          </h6>
        </Card.Header>
        <Card.Body className="scene-photo">
          <img
            alt=""
            src={getImage(scene.images, "landscape")}
            className="scene-photo-element"
          />
        </Card.Body>
        <Card.Footer className="row mx-1">
          <div className="scene-performers mr-auto">{performers}</div>
          {scene.duration && (
            <div title={`${scene.duration} seconds`}>
              Duration: <b>{formatDuration(scene.duration)}</b>
            </div>
          )}
          {scene.director && (
            <div className="ml-3">
              Director: <strong>{scene.director}</strong>
            </div>
          )}
        </Card.Footer>
      </Card>
      <Tabs
        activeKey={activeTab}
        id="scene-tabs"
        mountOnEnter
        onSelect={setTab}
      >
        <Tab eventKey="description" title="Description">
          <div className="scene-description my-4">
            <h4>Description:</h4>
            <div>{scene.details}</div>
            <div className="scene-tags">
              <h6>Tags:</h6>
              <ul className="scene-tag-list">{tags}</ul>
            </div>
            <hr />
            <div>
              <strong className="mr-2">Studio URL: </strong>
              <a
                href={getUrlByType(scene.urls, "STUDIO")}
                target="_blank"
                rel="noopener noreferrer"
              >
                {getUrlByType(scene.urls, "STUDIO")}
              </a>
            </div>
          </div>
        </Tab>
        <Tab eventKey="fingerprints" title="Fingerprints">
          <div className="scene-fingerprints my-4">
            <h4>Fingerprints:</h4>
            {fingerprints.length === 0 ? (
              <h6>No fingerprints found for this scene.</h6>
            ) : (
              <Table striped bordered size="sm">
                <thead>
                  <tr>
                    <td>
                      <b>Algorithm</b>
                    </td>
                    <td>
                      <b>Hash</b>
                    </td>
                    <td>
                      <b>Duration</b>
                    </td>
                    <td>
                      <b>Submissions</b>
                    </td>
                    <td>
                      <b>First Added</b>
                    </td>
                    <td>
                      <b>Last Added</b>
                    </td>
                  </tr>
                </thead>
                <tbody>{fingerprints}</tbody>
              </Table>
            )}
          </div>
        </Tab>
        <Tab
          eventKey="edits"
          title={`Edits${formatPendingEdits(pendingEditCount)}`}
          tabClassName={pendingEditCount ? "PendingEditTab" : ""}
        >
          <EditList type={TargetTypeEnum.SCENE} id={scene.id} />
        </Tab>
      </Tabs>
    </>
  );
};

export default SceneComponent;
