import Form from 'react-bootstrap/Form';

function SelectionFormat({ array, subject, onChange }) {
  const handleChange = (e) => {
    const value = e.target.value;
    onChange(value); // Pass the selected value to the parent
  };
  return (
    <>
      <Form.Label htmlFor={`create${subject}`}>Describe how much you like {subject} on a scale of 0-10</Form.Label>
      <Form.Select 
        aria-label={`${subject} select example`}
        onChange={handleChange}
      >
        <option>Describe how much you like {subject} on a scale of 0-10</option>
        {array.map((item, index) => (
          <option key={index} value={index}>
            {index} - {item}
          </option>
        ))}
      </Form.Select>
    </>
  );
}

export function HumanitySelectionFormat({ array, subject }) {
  return (
    <>
      <Form.Label htmlFor={`create${subject}`}>So it turns out {subject} doesn't exist on a scale - pick what makes sense to you.</Form.Label>
      <Form.Select aria-label={`${subject} select example`}>
        <option>Which option best describes what fits for you?</option>
        {array.map((item, index) => (
          <option key={index} value={index}>
            {item}
          </option>
        ))}
      </Form.Select>
    </>
  );
}


export default SelectionFormat;