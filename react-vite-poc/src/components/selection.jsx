import { useState } from 'react';
import Button from 'react-bootstrap/Button';
import Accordion from 'react-bootstrap/Accordion';

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
        <Accordion.Header>{subject} - {getHeaderText()}</Accordion.Header>
        <Accordion.Body>
          <div className="d-flex flex-column gap-2">
            {array.map((item, index) => (
              <Button
                key={index}
                variant={selected === index ? "primary" : "outline-secondary"}
                onClick={() => handleSelect(index)}
                className="text-start"
              >
                {index} - {item}
              </Button>
            ))}
          </div>
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

export default SelectionFormat;