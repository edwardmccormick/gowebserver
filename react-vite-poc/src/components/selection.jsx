import Form from 'react-bootstrap/Form';

function SelectionFormat({array, subject}) {

  return (
        <>     
        <Form.Label htmlFor="create{{subject}}">Name</Form.Label>
        <Form.Select aria-label="{{subect}} select example">
            <option>Open this {subject} menu</option>
                {array.map((item, index) => (
                <option key={index} value={index}>
                    {item}
                </option>
                ))}
        </Form.Select>
        </>
        )
}

export default SelectionFormat;