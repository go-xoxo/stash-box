import React from "react";
import { useHistory } from "react-router-dom";
import { Button, Col, Form } from "react-bootstrap";
import { useForm } from "react-hook-form";
import * as yup from "yup";
import { yupResolver } from "@hookform/resolvers/yup";

import { FullPerformer_findPerformer as Performer } from "src/graphql/definitions/FullPerformer";
import { usePerformerEdit, OperationEnum } from "src/graphql";
import { EditNote } from "src/components/form";
import { editHref } from "src/utils";

const schema = yup.object({
  id: yup.string().required(),
  note: yup.string().required("An edit note is required."),
});
export type FormData = yup.Asserts<typeof schema>;

interface Props {
  performer: Performer;
}

const PerformerDelete: React.FC<Props> = ({ performer }) => {
  const history = useHistory();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormData>({
    resolver: yupResolver(schema),
    mode: "onBlur",
  });
  const [deletePerformerEdit, { loading: deleting }] = usePerformerEdit({
    onCompleted: (data) => {
      if (data.performerEdit.id) history.push(editHref(data.performerEdit));
    },
  });

  const handleDelete = (data: FormData) =>
    deletePerformerEdit({
      variables: {
        performerData: {
          edit: {
            operation: OperationEnum.DESTROY,
            id: data.id,
            comment: data.note,
          },
        },
      },
    });

  return (
    <>
      <Form
        className="PerformerDeleteForm"
        onSubmit={handleSubmit(handleDelete)}
      >
        <Form.Row>
          <h4>
            Delete performer <em>{performer.name}</em>
          </h4>
        </Form.Row>
        <Form.Control type="hidden" value={performer.id} {...register("id")} />
        <Form.Row className="my-4">
          <Col md={6}>
            <EditNote register={register} error={errors.note} />
          </Col>
        </Form.Row>
        <Form.Row className="mt-2">
          <Button
            variant="danger"
            className="ml-auto mr-2"
            onClick={() => history.goBack()}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            disabled
            className="d-none"
            aria-hidden="true"
          />
          <Button type="submit" disabled={deleting}>
            Submit Edit
          </Button>
        </Form.Row>
      </Form>
    </>
  );
};

export default PerformerDelete;
