import React from "react";
import { useHistory } from "react-router-dom";
import { Button, Col, Form } from "react-bootstrap";
import { useForm } from "react-hook-form";
import * as yup from "yup";
import { yupResolver } from "@hookform/resolvers/yup";

import { useTagEdit, OperationEnum } from "src/graphql";
import { Tag_findTag as Tag } from "src/graphql/definitions/Tag";
import { EditNote } from "src/components/form";
import { editHref } from "src/utils";

const schema = yup.object({
  id: yup.string().required(),
  note: yup.string().required("An edit note is required."),
});
export type FormData = yup.Asserts<typeof schema>;

interface Props {
  tag: Tag;
}

const TagDelete: React.FC<Props> = ({ tag }) => {
  const history = useHistory();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormData>({
    resolver: yupResolver(schema),
    mode: "onBlur",
  });
  const [deleteTagEdit, { loading: deleting }] = useTagEdit({
    onCompleted: (data) => {
      if (data.tagEdit.id) history.push(editHref(data.tagEdit));
    },
  });

  const handleDelete = (data: FormData) =>
    deleteTagEdit({
      variables: {
        tagData: {
          edit: {
            operation: OperationEnum.DESTROY,
            id: data.id,
            comment: data.note,
          },
        },
      },
    });

  return (
    <Form className="TagDeleteForm" onSubmit={handleSubmit(handleDelete)}>
      <Form.Row>
        <h4>
          Delete tag <em>{tag.name}</em>
        </h4>
      </Form.Row>
      <Form.Control type="hidden" value={tag.id} {...register("id")} />
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
        <Button type="submit" disabled className="d-none" aria-hidden="true" />
        <Button type="submit" disabled={deleting}>
          Submit Edit
        </Button>
      </Form.Row>
    </Form>
  );
};

export default TagDelete;
