import SelectionFormat from './selection';
import details from '../../../details.json';

function DetailsSelections({ onChange}) {
  return (
    <div>
    <h4 className="text-center">I love it/I hate it:</h4>
        <div className='d-flex flex-column flex-sm-row flex-wrap justify-content-around'>
        {Object.keys(details).map((key) => (
            <div key={key} className="m-3 col-12 col-sm-5 col-md-3">
            <SelectionFormat 
              array={details[key]} 
              subject={key} 
              onChange={(value) => onChange(key, value)}
            />
            </div>
        ))}
        </div>
    </div>
  );
}

export default DetailsSelections;