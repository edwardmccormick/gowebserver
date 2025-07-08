import SelectionFormat from './selection';
import details from '../../../details.json';

function DetailsSelections({ onChange, selectedValues = {} }) {
  return (
    <div className="my-4">
      <h4 className="text-center mb-3">I love it/I hate it:</h4>
      <div className="d-flex flex-column">
        {Object.keys(details).map((key) => (
          <div key={key} className="mb-2">
            <SelectionFormat 
              array={details[key]} 
              subject={key} 
              onChange={(value) => onChange(key, value)}
              selectedValue={selectedValues[key] !== undefined ? selectedValues[key] : null}
            />
          </div>
        ))}
      </div>
    </div>
  );
}

export default DetailsSelections;