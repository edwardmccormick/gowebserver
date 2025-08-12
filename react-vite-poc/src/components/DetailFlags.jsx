import React from 'react';
import DetailFlag from './DetailFlag';
import { getNonEmptyDetails } from '../utils/formatters';

/**
 * Renders flags for a person's details
 * @param {Object} details - The details object containing scores
 * @param {number} limit - Max number of flags to show (0 for all)
 * @param {boolean} topValues - If true, show highest values first
 */
const DetailFlags = ({ details, limit = 0, topValues = false }) => {
  if (!details) {
    return null;
  }
  
  const nonEmptyDetails = getNonEmptyDetails(details);
  const detailKeys = Object.keys(nonEmptyDetails);
  
  // If there are no details with values, return null
  if (detailKeys.length === 0) {
    return null;
  }

  let keysToShow = detailKeys;
  
  // If topValues is true, sort by score (highest first)
  if (topValues) {
    keysToShow = [...detailKeys].sort((a, b) => {
      return nonEmptyDetails[b] - nonEmptyDetails[a];
    });
  }
  
  // If limit is specified and greater than 0, only show that many flags
  if (limit > 0) {
    keysToShow = keysToShow.slice(0, limit);
  }

  return (
    <div className="detail-flags">
      {keysToShow.map(key => (
        <DetailFlag key={key} detailKey={key} score={details[key]} />
      ))}
      
      {/* Show how many more details are available if limited */}
      {limit > 0 && detailKeys.length > limit && (
        <span className="badge bg-secondary">+{detailKeys.length - limit} more</span>
      )}
    </div>
  );
};

export default DetailFlags;