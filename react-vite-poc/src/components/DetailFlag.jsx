import React from 'react';
import { formatDetailKey } from '../utils/formatters';
import OverlayTrigger from 'react-bootstrap/OverlayTrigger';
import Tooltip from 'react-bootstrap/Tooltip';
import details from '../../../details.json';

// Map of detail keys to display names
const detailDisplayNames = {
  dogs: 'Dogs',
  cats: 'Cats',
  kids: 'Kids',
  smoking: 'Smoking',
  drinking: 'Drinking',
  religion: 'Religion',
  food: 'Food',
  energy_levels: 'Energy',
  outdoorsy_ness: 'Outdoorsy',
  travel: 'Travel',
  bougieness: 'Bougie',
  importance_of_politics: 'Politics'
};

const DetailFlag = ({ detailKey, score }) => {
  if (score === undefined || score === null) {
    return null;
  }

  // Normalize score to be between 0 and 10
  const normalizedScore = Math.max(0, Math.min(10, score));

  // Generate gradient color from red (low) through yellow to green (high)
  const getGradientColor = (score) => {
    // Red component starts at 255 and remains 255 until score 5, then decreases to 0 as score approaches 10
    const red = score <= 5 ? 255 : Math.round(255 * (10 - score) / 5);
    
    // Green component increases from 0 to 255 as score approaches 5, then stays 255 until 10
    const green = score <= 5 ? Math.round(255 * score / 5) : 255;
    
    // Blue is always 0 for this red->yellow->green gradient
    return `rgb(${red}, ${green}, 0)`;
  };

  const flagColor = getGradientColor(normalizedScore);
  
  // Get display name for the detail
  const displayName = detailDisplayNames[detailKey] || formatDetailKey(detailKey);
  
  // Get the description text for the tooltip from details.json
  let tooltipText = "No description available";
  if (details[detailKey] && Array.isArray(details[detailKey]) && details[detailKey][normalizedScore]) {
    tooltipText = details[detailKey][normalizedScore];
  }

  // Create the flag element with tooltip
  const flagElement = (
    <span className="badge text-white" style={{ backgroundColor: flagColor, fontSize: '0.85rem', cursor: 'pointer' }}>
      <svg 
        xmlns="http://www.w3.org/2000/svg" 
        width="16" 
        height="16" 
        fill="currentColor" 
        className="bi bi-flag-fill me-1" 
        viewBox="0 0 16 16"
      >
        <path d="M14.778.085A.5.5 0 0 1 15 .5V8a.5.5 0 0 1-.314.464L14.5 8l.186.464-.003.001-.006.003-.023.009a12.435 12.435 0 0 1-.397.15c-.264.095-.631.223-1.047.35-.816.252-1.879.523-2.71.523-.847 0-1.548-.28-2.158-.525l-.028-.01C7.68 8.71 7.14 8.5 6.5 8.5c-.7 0-1.638.23-2.437.477A19.626 19.626 0 0 0 3 9.342V15.5a.5.5 0 0 1-1 0V.5a.5.5 0 0 1 1 0v.282c.226-.079.496-.17.79-.26C4.606.272 5.67 0 6.5 0c.84 0 1.524.277 2.121.519l.043.018C9.286.788 9.828 1 10.5 1c.7 0 1.638-.23 2.437-.477a19.587 19.587 0 0 0 1.349-.476l.019-.007.004-.002h.001"/>
      </svg>
      {displayName}: {normalizedScore}
    </span>
  );

  return (
    <div className="detail-flag d-inline-flex align-items-center me-2 mb-2">
      <OverlayTrigger
        placement="top"
        delay={{ show: 100, hide: 100 }}
        overlay={(props) => (
          <Tooltip id={`tooltip-${detailKey}-${normalizedScore}`} {...props}>
            {tooltipText}
          </Tooltip>
        )}
      >
        {flagElement}
      </OverlayTrigger>
    </div>
  );
};

export default DetailFlag;