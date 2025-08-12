/**
 * Format a detail key into a display-friendly string
 * For example: "outdoorsy_ness" -> "Outdoorsy Ness"
 * 
 * @param {string} key - The detail key to format
 * @returns {string} Formatted string
 */
export const formatDetailKey = (key) => {
  if (!key) return '';
  
  return key
    .replace(/_/g, ' ')              // Replace underscores with spaces
    .split(' ')                      // Split into words
    .map(word => word.charAt(0).toUpperCase() + word.slice(1)) // Capitalize first letter of each word
    .join(' ');                      // Join back with spaces
};

/**
 * Get only details that have values
 * 
 * @param {Object} details - The details object
 * @returns {Object} Object with only defined values
 */
export const getNonEmptyDetails = (details) => {
  if (!details) return {};
  
  return Object.entries(details)
    .filter(([_, value]) => value !== null && value !== undefined)
    .reduce((acc, [key, value]) => {
      acc[key] = value;
      return acc;
    }, {});
};