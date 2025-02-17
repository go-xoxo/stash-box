import React from "react";
import { Button, OverlayTrigger, Popover } from "react-bootstrap";
import { Icon } from "src/components/fragments";
import { faQuestionCircle } from "@fortawesome/free-solid-svg-icons";

interface Props {
  message: string;
}

const Help: React.FC<Props> = ({ message }) => {
  const renderContent = () => (
    <Popover id="help">
      <Popover.Content>{message}</Popover.Content>
    </Popover>
  );

  return (
    <OverlayTrigger
      overlay={renderContent()}
      placement="bottom"
      trigger="hover"
    >
      <Button variant="link" className="minimal">
        <Icon icon={faQuestionCircle} />
      </Button>
    </OverlayTrigger>
  );
};

export default Help;
