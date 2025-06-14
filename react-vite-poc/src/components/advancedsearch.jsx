import React, { useState } from 'react';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import InputGroup from 'react-bootstrap/InputGroup';
import humanity from '../../../humanity.json';
import details from '../../../details.json';

function AdvancedSearch({ onSearch }) {
  const [distance, setDistance] = useState(50); // Default distance
  const [criteria, setCriteria] = useState({});
  const [selectedCategory, setSelectedCategory] = useState('');
  const [selectedValue, setSelectedValue] = useState('');

  const handleAddCriteria = () => {
    if (selectedCategory && selectedValue) {
      setCriteria((prev) => ({
        ...prev,
        [selectedCategory]: selectedValue,
      }));
      setSelectedCategory('');
      setSelectedValue('');
    }
  };

  const handleSearch = () => {
    onSearch({ distance, criteria });
  };

  return (
    <div className="p-4 border rounded shadow-sm bg-light">
      <h3 className="text-center mb-4">Advanced Search</h3>
      <Form>
        {/* Distance Slider */}
        <Form.Group className="mb-3">
          <Form.Label>Maximum Distance: {distance} miles</Form.Label>
          <Form.Range
            min={1}
            max={500}
            value={distance}
            onChange={(e) => setDistance(e.target.value)}
          />
        </Form.Group>

        {/* Criteria Selection */}
        <Form.Group className="mb-3">
          <Form.Label>Select Criteria</Form.Label>
          <InputGroup>
            <Form.Select
              value={selectedCategory}
              onChange={(e) => setSelectedCategory(e.target.value)}
            >
              <option value="">Select Category</option>
              {Object.keys({ ...humanity, ...details }).map((key) => (
                <option key={key} value={key}>
                  {key.replace(/_/g, ' ')}
                </option>
              ))}
            </Form.Select>
            <Form.Select
              value={selectedValue}
              onChange={(e) => setSelectedValue(e.target.value)}
              disabled={!selectedCategory}
            >
              <option value="">Select Value</option>
              {selectedCategory &&
                (humanity[selectedCategory] || details[selectedCategory]).map(
                  (item, index) => (
                    <option key={index} value={index}>
                      {item}
                    </option>
                  )
                )}
            </Form.Select>
            <Button variant="primary" onClick={handleAddCriteria}>
              Add
            </Button>
          </InputGroup>
        </Form.Group>

        {/* Display Selected Criteria */}
        {Object.keys(criteria).length > 0 && (
          <div className="mb-3">
            <h5>Selected Criteria:</h5>
            <ul className="list-group">
              {Object.entries(criteria).map(([key, value]) => (
                <li key={key} className="list-group-item">
                  <strong>{key.replace(/_/g, ' ')}:</strong>{' '}
                  {(humanity[key] || details[key])[value]}
                </li>
              ))}
            </ul>
          </div>
        )}

        {/* Search Button */}
        <Button
          variant="success"
          className="w-100 mt-3"
          onClick={handleSearch}
        >
          Search
        </Button>
      </Form>
    </div>
  );
}

export default AdvancedSearch;