import Form from 'react-bootstrap/Form';

function SelectionFormat({ array, subject }) {
  return (
    <>
      <Form.Label htmlFor={`create${subject}`}>Describe how much you like {subject} on a scale of 0-10</Form.Label>
      <Form.Select aria-label={`${subject} select example`}>
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

export function NonLinearSelectionFormat({ array, subject }) {
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