import React from "react";
import { Link } from "react-router-dom";
import { Card } from "react-bootstrap";
import { faVideo } from "@fortawesome/free-solid-svg-icons";

import { Scenes_queryScenes_scenes as Performance } from "src/graphql/definitions/Scenes";
import { getImage, sceneHref, studioHref, formatDuration } from "src/utils";
import { Icon } from "src/components/fragments";

const CLASSNAME = "SceneCard";
const CLASSNAME_IMAGE = `${CLASSNAME}-image`;
const CLASSNAME_TITLE = `${CLASSNAME}-title`;
const CLASSNAME_BODY = `${CLASSNAME}-body`;

const SceneCard: React.FC<{ performance: Performance }> = ({ performance }) => (
  <div className={`col-3 ${CLASSNAME}`}>
    <Card>
      <Card.Body className={CLASSNAME_BODY}>
        <Link to={sceneHref(performance)} className={CLASSNAME_IMAGE}>
          <img alt="" src={getImage(performance.images, "landscape")} />
        </Link>
      </Card.Body>
      <Card.Footer>
        <Link to={sceneHref(performance)} className="d-flex">
          <h6 className={CLASSNAME_TITLE}>{performance.title}</h6>
          <span className="text-muted">
            {performance.duration ? formatDuration(performance.duration) : ""}
          </span>
        </Link>
        <div className="text-muted">
          {performance.studio && (
            <Link
              to={studioHref(performance.studio)}
              className="float-right text-truncate SceneCard-studio-name"
            >
              <Icon icon={faVideo} className="mr-1" />
              {performance.studio.name}
            </Link>
          )}
          <strong>{performance.date}</strong>
        </div>
      </Card.Footer>
    </Card>
  </div>
);

export default SceneCard;
