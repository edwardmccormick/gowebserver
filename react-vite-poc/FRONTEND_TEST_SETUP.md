# Frontend Testing Setup Guide

This guide provides detailed instructions for setting up and running tests for the React frontend of our dating web application.

## Table of Contents

1. [Installation](#installation)
2. [Running Tests](#running-tests)
3. [Test Coverage](#test-coverage)
4. [Writing Tests](#writing-tests)
5. [Mocking Strategies](#mocking-strategies)
6. [Common Issues](#common-issues)

## Installation

First, make sure all required testing dependencies are installed:

```bash
cd react-vite-poc
npm install --save-dev jest @testing-library/react @testing-library/user-event @testing-library/jest-dom jest-environment-jsdom @babel/preset-env @babel/preset-react babel-jest
```

The following packages should be listed in your `package.json` devDependencies:

```json
"devDependencies": {
  "@babel/preset-env": "^7.x.x",
  "@babel/preset-react": "^7.x.x",
  "@testing-library/jest-dom": "^5.x.x",
  "@testing-library/react": "^14.x.x",
  "@testing-library/user-event": "^14.x.x",
  "babel-jest": "^29.x.x",
  "jest": "^29.x.x",
  "jest-environment-jsdom": "^29.x.x"
}
```

Create a `.babelrc` file in the project root:

```json
{
  "presets": [
    ["@babel/preset-env", { "targets": { "node": "current" } }],
    ["@babel/preset-react", { "runtime": "automatic" }]
  ]
}
```

## Running Tests

Add the following script to your `package.json`:

```json
"scripts": {
  "test": "jest",
  "test:watch": "jest --watch",
  "test:coverage": "jest --coverage"
}
```

To run all tests:

```bash
npm test
```

To run tests in watch mode (recommended during development):

```bash
npm run test:watch
```

To run specific test files:

```bash
npm test -- src/components/__tests__/NotificationService.test.jsx
```

## Test Coverage

View test coverage report after running:

```bash
npm run test:coverage
```

The coverage report will be generated in the `coverage` directory. Open `coverage/lcov-report/index.html` in a browser to view a detailed report.

Target coverage thresholds are defined in `jest.config.js`:

```javascript
coverageThreshold: {
  global: {
    statements: 50,
    branches: 40,
    functions: 50,
    lines: 50,
  },
}
```

## Writing Tests

### Component Testing Pattern

Follow this structure for component tests:

```jsx
import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import YourComponent from '../your-component';

// Optional: Mock child components to focus on the component being tested
jest.mock('../child-component', () => {
  return function MockChildComponent(props) {
    return <div data-testid="mock-child">{props.children}</div>;
  };
});

describe('YourComponent', () => {
  // Setup common test data
  const mockProps = {
    // Component props here
  };
  
  test('renders correctly', () => {
    render(<YourComponent {...mockProps} />);
    expect(screen.getByText('Expected Text')).toBeInTheDocument();
  });
  
  test('handles user interaction', () => {
    render(<YourComponent {...mockProps} />);
    const button = screen.getByRole('button', { name: /click me/i });
    fireEvent.click(button);
    expect(screen.getByText('After Click')).toBeInTheDocument();
  });
});
```

### Testing Real-time Features

For WebSockets and SSE:

1. Mock the browser APIs (EventSource, WebSocket)
2. Use the provided mock implementations in tests
3. Manually trigger events to simulate server messages

Example for SSE testing:

```jsx
// See NotificationService.test.jsx for a complete example
```

## Mocking Strategies

### API Mocking

For API calls, mock the fetch function:

```jsx
// In your test file
beforeEach(() => {
  jest.spyOn(global, 'fetch').mockImplementation((url) => {
    if (url.includes('/api/users')) {
      return Promise.resolve({
        ok: true,
        json: () => Promise.resolve([{ id: 1, name: 'Test User' }])
      });
    }
    return Promise.reject(new Error('Not found'));
  });
});

afterEach(() => {
  global.fetch.mockRestore();
});
```

### Component Mocking

Mock child components when you want to focus on testing the parent component's behavior:

```jsx
jest.mock('../ChatModal', () => {
  return function MockChatModal(props) {
    return (
      <div data-testid="chat-modal">
        Mock Chat Modal for {props.person?.name}
      </div>
    );
  };
});
```

### Browser API Mocking

For browser APIs like localStorage or EventSource:

```jsx
// See setupTests.js for implementation examples
```

## Common Issues

### Module Import Errors

If you encounter errors like "Cannot find module", update the moduleNameMapper in jest.config.js:

```js
moduleNameMapper: {
  '^@/(.*)$': '<rootDir>/src/$1',
}
```

### React Version Errors

If you're using React 18 or newer, make sure your testing library versions are compatible:

- Use @testing-library/react version 13.0.0 or newer
- Use @testing-library/user-event version 14.0.0 or newer

### CSS/Asset Import Errors

CSS and asset imports are mocked in the test environment:

- CSS files are mocked in `src/__mocks__/styleMock.js`
- Image/asset files are mocked in `src/__mocks__/fileMock.js`

If you encounter new file types that cause issues, add them to the moduleNameMapper in jest.config.js.