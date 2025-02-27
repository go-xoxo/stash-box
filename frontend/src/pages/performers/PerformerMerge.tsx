import React, { useState } from "react";
import { useHistory } from "react-router-dom";
import { Button, Col, Form, Row } from "react-bootstrap";
import { flatMap } from "lodash-es";

import { FullPerformer_findPerformer as Performer } from "src/graphql/definitions/FullPerformer";
import { SearchPerformers_searchPerformer as SearchPerformer } from "src/graphql/definitions/SearchPerformers";
import {
  usePerformerEdit,
  OperationEnum,
  PerformerEditDetailsInput,
} from "src/graphql";

import PerformerSelect from "src/components/performerSelect";
import PerformerCard from "src/components/performerCard";
import { editHref } from "src/utils";
import PerformerForm from "./performerForm";

const CLASSNAME = "PerformerMerge";

interface Props {
  performer: Performer;
}

const PerformerMerge: React.FC<Props> = ({ performer }) => {
  const history = useHistory();
  const [mergeActive, setMergeActive] = useState(false);
  const [mergeSources, setMergeSources] = useState<SearchPerformer[]>([]);
  const [aliasUpdating, setAliasUpdating] = useState(true);
  const [insertPerformerEdit, { loading: saving }] = usePerformerEdit({
    onCompleted: (data) => {
      if (data.performerEdit.id) history.push(editHref(data.performerEdit));
    },
  });

  const doUpdate = (
    insertData: PerformerEditDetailsInput,
    editNote: string,
    setModifyAliases: boolean
  ) => {
    insertPerformerEdit({
      variables: {
        performerData: {
          edit: {
            id: performer.id,
            operation: OperationEnum.MERGE,
            merge_source_ids: mergeSources.map((p) => p.id),
            comment: editNote,
          },
          details: insertData,
          options: {
            set_merge_aliases: aliasUpdating,
            set_modify_aliases: setModifyAliases,
          },
        },
      },
    });
  };

  return (
    <div className={CLASSNAME}>
      <h3>
        Merge performers into <em>{performer.name}</em>
      </h3>
      <hr />
      <div className="row">
        <div className="col-6">
          {!mergeActive && (
            <>
              <PerformerSelect
                performers={[]}
                onChange={(performers) => setMergeSources(performers)}
                message="Search for performers to merge..."
                excludePerformers={[
                  performer.id,
                  ...mergeSources.map((p) => p.id),
                ]}
              />
              {mergeSources.length > 0 && (
                <Button
                  onClick={() => setMergeActive(true)}
                  className="ml-auto"
                >
                  Continue
                </Button>
              )}
            </>
          )}
          {mergeActive && (
            <Row>
              <Col xs={3}>
                <h6 className="text-center">Merge Target</h6>
                <PerformerCard performer={performer} className="TargetCard" />
              </Col>
              <Col xs={9}>
                <Row className="mt-4">
                  {mergeSources.map((source) => (
                    <Col xs={4} key={source.id}>
                      <PerformerCard performer={source} />
                    </Col>
                  ))}
                </Row>
              </Col>
            </Row>
          )}
        </div>
        <div className="col-6">
          <p>
            Merging performers reassigns all scene performances of the sources
            to the target performer. The source <i>stashIds</i> will be
            redirected so that previously tagged content can be updated with the
            new performers metadata.
          </p>
          <p>
            This operation is not easily reversible and attention should be paid
            that all entities are truly the same.
          </p>
        </div>
      </div>
      {mergeActive && (
        <>
          <Form.Check
            id="merge-alias-updating"
            checked={aliasUpdating}
            onChange={() => setAliasUpdating(!aliasUpdating)}
            label="Update scene performance aliases on merged performers to old performer name."
          />
          <h5 className="mt-4">
            Update performer metadata for <em>{performer.name}</em>
          </h5>
          <PerformerForm
            performer={performer}
            initialAliases={[
              ...mergeSources.map((p) => p.name),
              ...flatMap(mergeSources, (p) => p.aliases),
            ]}
            initialImages={flatMap(mergeSources, (i) => i.images)}
            callback={doUpdate}
            changeType="merge"
            saving={saving}
          />
        </>
      )}
    </div>
  );
};

export default PerformerMerge;
