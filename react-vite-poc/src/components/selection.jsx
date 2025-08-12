import { useState } from 'react';
import Button from 'react-bootstrap/Button';
import Accordion from 'react-bootstrap/Accordion';
import DetailFlag from './DetailFlag';
import OverlayTrigger from 'react-bootstrap/OverlayTrigger';
import Tooltip from 'react-bootstrap/Tooltip';
import details from '../../../details.json';

function SelectionFormat({ array, subject, onChange, selectedValue = null }) {
  const [selected, setSelected] = useState(selectedValue);

  const handleSelect = (index) => {
    setSelected(index);
    onChange(index);
  };
  
  // Function to determine what to display in the header
  const getHeaderText = () => {
    return selected !== null && selected !== undefined ? array[selected] : `Love 'em or hate 'em?`;
  };

  return (
    <Accordion className="mb-3">
      <Accordion.Item eventKey="0">
        <Accordion.Header>
          <div className="d-flex align-items-center">
            <div className="me-2">{subject}</div>
            {selected !== null && selected !== undefined && (
              <DetailFlag detailKey={subject} score={selected} />
            )}
            {(selected === null || selected === undefined) && (
              <span className="text-secondary ms-2">- {getHeaderText()}</span>
            )}
          </div>
        </Accordion.Header>
        <Accordion.Body>
          <div className="d-flex flex-column gap-2">
            {array.map((item, index) => (
              <Button
                key={index}
                variant={selected === index ? "primary" : "outline-secondary"}
                onClick={() => handleSelect(index)}
                className="text-start d-flex justify-content-between align-items-center"
              >
                <span>{index} - {item}</span>
                {selected === index && (
                  <OverlayTrigger
                    placement="right"
                    delay={{ show: 100, hide: 100 }}
                    overlay={(props) => (
                      <Tooltip id={`tooltip-button-${subject}-${index}`} {...props}>
                        {item}
                      </Tooltip>
                    )}
                  >
                    <div className="d-inline-block">
                      <DetailFlag detailKey={subject} score={index} />
                    </div>
                  </OverlayTrigger>
                )}
              </Button>
            ))}
          </div>
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

export default SelectionFormat;