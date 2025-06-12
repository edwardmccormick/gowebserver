import SelectionFormat from './selection';
import verbiage from '../../../verbiage.json';

function VerbiageSelections() {
  return (
    <div>
    <h4 className="text-center">I love it/I hate it:</h4>
        <div className='d-flex flex-column flex-sm-row flex-wrap justify-content-around'>
        {Object.keys(verbiage).map((key) => (
            <div key={key} className="m-3 col-12 col-sm-5 col-md-3">
            <SelectionFormat array={verbiage[key]} subject={key} />
            </div>
        ))}
        </div>
    </div>
  );
}

export default VerbiageSelections;