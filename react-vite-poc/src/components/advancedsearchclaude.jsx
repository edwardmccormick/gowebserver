import React, { useState } from 'react';
import { Search, Plus, X, MapPin, Filter, Heart, Users } from 'lucide-react';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import Card from 'react-bootstrap/Card';
import Badge from 'react-bootstrap/Badge';
import InputGroup from 'react-bootstrap/InputGroup';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import details from '../../../details.json';

function ClaudeAdvancedSearch({ onSearch }) {
  const [distance, setDistance] = useState(50);
  const [gender, setGender] = useState('?')
  const [preference, setPreference] = useState('?')
  const [relationship, setRelationship] = useState('?');
  const [criteria, setCriteria] = useState({});
  const [selectedCategory, setSelectedCategory] = useState('');
  const [selectedValue, setSelectedValue] = useState('');
  const [searchType, setSearchType] = useState('exact'); // 'exact', 'min', 'max', 'range'
  const [rangeMin, setRangeMin] = useState('');
  const [rangeMax, setRangeMax] = useState('');
  const [isExpanded, setIsExpanded] = useState(false);

  const handleAddCriteria = () => {
    if (selectedCategory && (selectedValue !== '' || (searchType === 'range' && rangeMin !== '' && rangeMax !== ''))) {
      const criteriaKey = `${selectedCategory}_${Date.now()}`;
      let criteriaValue;
      
      switch (searchType) {
        case 'exact':
          criteriaValue = { type: 'exact', value: parseInt(selectedValue) };
          break;
        case 'min':
          criteriaValue = { type: 'min', value: parseInt(selectedValue) };
          break;
        case 'max':
          criteriaValue = { type: 'max', value: parseInt(selectedValue) };
          break;
        case 'range':
          criteriaValue = { type: 'range', min: parseInt(rangeMin), max: parseInt(rangeMax) };
          break;
        default:
          criteriaValue = { type: 'exact', value: parseInt(selectedValue) };
      }
      
      setCriteria((prev) => ({
        ...prev,
        [criteriaKey]: {
          category: selectedCategory,
          ...criteriaValue,
          display: getCriteriaDisplay(selectedCategory, criteriaValue)
        },
      }));
      
      // Reset form
      setSelectedCategory('');
      setSelectedValue('');
      setRangeMin('');
      setRangeMax('');
      setSearchType('exact');
      setGender('?');
      setPreference('?');
    }
  };

  const getCriteriaDisplay = (category, criteriaValue) => {
    const categoryItems = details[category] || [];
    const categoryName = category.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase());
    
    switch (criteriaValue.type) {
      case 'exact':
        return `${categoryName}: ${criteriaValue.value} exactly - ${categoryItems[criteriaValue.value] || criteriaValue.value}`;
      case 'min':
        return `${categoryName}: ${criteriaValue.value} or higher - ${categoryItems[criteriaValue.value] || criteriaValue.value} or higher`;
      case 'max':
        return `${categoryName}: ${criteriaValue.value} or lower - ${categoryItems[criteriaValue.value] || criteriaValue.value} or lower`;
      case 'range':
        return `${categoryName}: Between ${criteriaValue.min} and ${criteriaValue.max} - ${categoryItems[criteriaValue.min] || criteriaValue.min} and ${categoryItems[criteriaValue.max] || criteriaValue.max}`;
      default:
        return `${categoryName}: ${criteriaValue.value}`;

    // switch (criteriaValue.type) {
    //   case 'exact':
    //     return `${categoryName}: ${criteriaValue.value}`;
    //   case 'min':
    //     return `${categoryName}: ${criteriaValue.value}`;
    //   case 'max':
    //     return `${categoryName}: ${criteriaValue.value}`;
    //   case 'range':
    //     return `${categoryName}: ${criteriaValue.max}`;
    //   default:
    //     return `${categoryName}: ${criteriaValue.value}`;

    }
  };

  const removeCriteria = (key) => {
    setCriteria((prev) => {
      const newCriteria = { ...prev };
      delete newCriteria[key];
      return newCriteria;
    });
  };

  const handleSearch = () => {
    onSearch({ distance, criteria });
  };

  const clearAllCriteria = () => {
    setCriteria({});
  };

  const getDistanceColor = () => {
    const percentage = (distance / 500) * 100;
    if (percentage < 20) return 'success';
    if (percentage < 40) return 'info';
    if (percentage < 60) return 'warning';
    if (percentage < 80) return 'danger';
    return 'dark';
  };

  return (
    <div className="container-fluid mb-4">
      <style jsx>{`
        .gradient-header {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          color: white;
        }
        .filter-card {
          border: none;
          box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
          transition: all 0.3s ease;
        }
        .filter-card:hover {
          box-shadow: 0 8px 15px rgba(0, 0, 0, 0.15);
          transform: translateY(-2px);
        }
        .distance-slider {
          background: linear-gradient(to right, #28a745, #ffc107, #dc3545);
          height: 8px;
          border-radius: 4px;
        }
        .search-btn {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          border: none;
          transition: all 0.3s ease;
        }
        .search-btn:hover {
          transform: translateY(-2px);
          box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
        }
        .criteria-badge {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          color: white;
          border: none;
          cursor: pointer;
          transition: all 0.2s ease;
        }
        .criteria-badge:hover {
          transform: scale(1.05);
          box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
        }
        .category-label {
          font-weight: 600;
          color: #495057;
          margin-bottom: 0.5rem;
        }
        .expand-btn {
          background: rgba(255, 255, 255, 0.2);
          border: 1px solid rgba(255, 255, 255, 0.3);
          color: white;
          transition: all 0.2s ease;
        }
        .expand-btn:hover {
          background: rgba(255, 255, 255, 0.3);
          color: white;
        }
      `}</style>

      <Card className="filter-card">
        {/* Header */}
        <Card.Header className="gradient-header">
          <Row className="align-items-center">
            <Col>
              <div className="d-flex align-items-center">
                <div className="me-3 p-2 bg-white bg-opacity-25 rounded">
                  <Search size={24} />
                </div>
                <div>
                  <h3 className="mb-1">Find Someone Awesome</h3>
                  <small className="opacity-75">Customize your search preferences</small>
                </div>
              </div>
            </Col>
            <Col xs="auto">
              <Button
                variant="outline-light"
                className="expand-btn"
                onClick={() => setIsExpanded(!isExpanded)}
              >
                <Filter size={16} />
              </Button>
            </Col>
          </Row>
        </Card.Header>

        <Card.Body>
          {/* Distance Section */}
          <div className="mb-4 pb-4 border-bottom">
            <div className="d-flex align-items-center mb-3">
              <MapPin size={20} className="text-primary me-2" />
              <h5 className="mb-0">Distance: {distance} miles - I am {gender} interested in {preference} for a {relationship}</h5>
            </div>
            
            {/* <Form.Range
              min={1}
              max={500}
              value={distance}
              onChange={(e) => setDistance(e.target.value)}
              className="mb-2"
            />
            
            <div className="d-flex justify-content-between">
              <small className="text-muted">1 mile</small>
              <Badge bg={getDistanceColor()}>{distance} miles</Badge>
              <small className="text-muted">500 miles</small>
            </div> */}
            <div className="d-flex flex-row justify-content-between align-items-center">
                <Form.Select
                value={distance}
                onChange={(e) => setDistance(e.target.value)}
                className="mb-2 col-3">
                    <option value="1">How far would you go for love?</option>
                    <option value="10">10 miles</option>
                    <option value="25">25 miles</option>
                    <option value="50">50 miles</option>
                    <option value="100">100 miles</option>
                    <option value="250">250 miles</option>
                    <option value="500">500 miles</option>
                </Form.Select>
                <Form.Select
                value={gender}
                onChange={(e) => setGender(e.target.value)}
                className="mb-2 col-3">
                    <option value="?">I am....</option>
                    <option value="a woman">a woman</option>
                    <option value="a man">a man</option>
                    <option value="nonbinary">nonbinary</option>
                    <option value="who cares">does it matter?</option>
                </Form.Select>
                <Form.Select
                value={preference}
                onChange={(e) => setPreference(e.target.value)}
                className="mb-2 col-3">
                    <option value="?">I am interested in....</option>
                    <option value="woman">women</option>
                    <option value="man">men</option>
                    <option value="nonbinary folks">nonbinary folks</option>
                    <option value="anybody">does it matter?</option>
                </Form.Select>
                <Form.Select
                value={relationship}
                onChange={(e) => setRelationship(e.target.value)}
                className="mb-2 col-3">
                    <option value="?">I want....</option>
                    <option value="something casual">something casual</option>
                    <option value="dating">dating</option>
                    <option value="something serious">something serious</option>
                    <option value="something kinky">something kinky</option>
                    <option value="who cares">does it matter?</option>
                </Form.Select>
            </div>
          </div>

          {/* Criteria Builder */}
          {(isExpanded || Object.keys(criteria).length > 0) && (
            <div className="mb-4 pb-4 border-bottom bg-light p-3 rounded">
              <div className="d-flex align-items-center mb-3">
                <Heart size={20} className="text-danger me-2" />
                <h5 className="mb-0">Preference Filters</h5>
              </div>
              
              <Row className="mb-3">
                <Col md={6}>
                  <Form.Label className="category-label">Category</Form.Label>
                  <Form.Select
                    value={selectedCategory}
                    onChange={(e) => setSelectedCategory(e.target.value)}
                    className="mb-2"
                  >
                    <option value="">Select a category...</option>
                    {Object.keys(details).map((key) => (
                      <option key={key} value={key}>
                        {key.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())}
                      </option>
                    ))}
                  </Form.Select>
                </Col>

                <Col md={6}>
                  <Form.Label className="category-label">Search Type</Form.Label>
                  <Form.Select
                    value={searchType}
                    onChange={(e) => setSearchType(e.target.value)}
                    className="mb-2"
                  >
                    <option value="exact">Exact Match</option>
                    <option value="min">Minimum Level</option>
                    <option value="max">Maximum Level</option>
                    <option value="range">Range</option>
                  </Form.Select>
                </Col>
              </Row>

              {/* Value Selection */}
              {selectedCategory && searchType !== 'range' && (
                <div className="mb-3">
                  <Form.Label className="category-label">
                    {searchType === 'exact' ? 'Select Value' : 
                     searchType === 'min' ? 'Minimum Value' : 'Maximum Value'}
                  </Form.Label>
                  <Form.Select
                    value={selectedValue}
                    onChange={(e) => setSelectedValue(e.target.value)}
                  >
                    <option value="">Choose an option...</option>
                    {(details[selectedCategory] || []).map((item, index) => (
                      <option key={index} value={index}>
                        {index} - {item.length > 80 ? `${item.substring(0, 80)}...` : item}
                      </option>
                    ))}
                  </Form.Select>
                </div>
              )}

              {/* Range Selection */}
              {selectedCategory && searchType === 'range' && (
                <Row className="mb-3">
                  <Col md={6}>
                    <Form.Label className="category-label">From</Form.Label>
                    <Form.Select
                      value={rangeMin}
                      onChange={(e) => setRangeMin(e.target.value)}
                    >
                      <option value="">Select minimum...</option>
                      {(details[selectedCategory] || []).map((item, index) => (
                        <option key={index} value={index}>
                          {item.length > 40 ? `${item.substring(0, 40)}...` : item}
                        </option>
                      ))}
                    </Form.Select>
                  </Col>
                  <Col md={6}>
                    <Form.Label className="category-label">To</Form.Label>
                    <Form.Select
                      value={rangeMax}
                      onChange={(e) => setRangeMax(e.target.value)}
                    >
                      <option value="">Select maximum...</option>
                      {(details[selectedCategory] || []).map((item, index) => (
                        <option key={index} value={index}>
                          {item.length > 40 ? `${item.substring(0, 40)}...` : item}
                        </option>
                      ))}
                    </Form.Select>
                  </Col>
                </Row>
              )}

              <Button
                onClick={handleAddCriteria}
                disabled={!selectedCategory || (searchType !== 'range' && selectedValue === '') || (searchType === 'range' && (rangeMin === '' || rangeMax === ''))}
                variant="primary"
                className="d-flex align-items-center"
              >
                <Plus size={16} className="me-2" />
                Add Filter
              </Button>
            </div>
          )}

          {/* Selected Criteria */}
          {Object.keys(criteria).length > 0 && (
            <div className="mb-4">
              <div className="d-flex justify-content-between align-items-center mb-3">
                <div className="d-flex align-items-center">
                  <Users size={20} className="text-info me-2" />
                  <h5 className="mb-0">Active Filters ({Object.keys(criteria).length})</h5>
                </div>
                <Button
                  variant="outline-danger"
                  size="sm"
                  onClick={clearAllCriteria}
                >
                  Clear All
                </Button>
              </div>
              
              <div className="d-flex flex-wrap gap-2">
                {Object.entries(criteria).map(([key, value]) => (
                  <Badge
                    key={key}
                    className="criteria-badge p-2 d-flex align-items-center"
                    style={{ fontSize: '0.85rem' }}
                  >
                    <span className="me-2">
                      {value.display.length > 60 ? `${value.display.substring(0, 60)}...` : value.display}
                    </span>
                    <Button
                      variant="link"
                      size="sm"
                      className="p-0 text-white"
                      onClick={() => removeCriteria(key)}
                      style={{ fontSize: '0.8rem', lineHeight: 1 }}
                    >
                      <X size={14} />
                    </Button>
                  </Badge>
                ))}
              </div>
            </div>
          )}

          {/* Search Button */}
          <Button
            onClick={handleSearch}
            className="search-btn w-100 py-3"
            size="lg"
          >
            <Search size={20} className="me-2" />
            Find My Matches
          </Button>
        </Card.Body>
      </Card>
    </div>
  );
}

export default ClaudeAdvancedSearch;