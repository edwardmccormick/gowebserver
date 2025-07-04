"use client";
import {
  InputGroupContext_default
} from "./chunk-DAZKUHWX.js";
import {
  FormCheckInput_default
} from "./chunk-LOC5JTWD.js";
import {
  require_classnames,
  require_jsx_runtime,
  useBootstrapPrefix
} from "./chunk-5TWWBDIN.js";
import {
  __toESM,
  require_react
} from "./chunk-5WQJO2FO.js";

// node_modules/react-bootstrap/esm/InputGroup.js
var import_classnames2 = __toESM(require_classnames());
var React2 = __toESM(require_react());
var import_react = __toESM(require_react());

// node_modules/react-bootstrap/esm/InputGroupText.js
var React = __toESM(require_react());
var import_classnames = __toESM(require_classnames());
var import_jsx_runtime = __toESM(require_jsx_runtime());
var InputGroupText = React.forwardRef(({
  className,
  bsPrefix,
  as: Component = "span",
  ...props
}, ref) => {
  bsPrefix = useBootstrapPrefix(bsPrefix, "input-group-text");
  return (0, import_jsx_runtime.jsx)(Component, {
    ref,
    className: (0, import_classnames.default)(className, bsPrefix),
    ...props
  });
});
InputGroupText.displayName = "InputGroupText";
var InputGroupText_default = InputGroupText;

// node_modules/react-bootstrap/esm/InputGroup.js
var import_jsx_runtime2 = __toESM(require_jsx_runtime());
var InputGroupCheckbox = (props) => (0, import_jsx_runtime2.jsx)(InputGroupText_default, {
  children: (0, import_jsx_runtime2.jsx)(FormCheckInput_default, {
    type: "checkbox",
    ...props
  })
});
var InputGroupRadio = (props) => (0, import_jsx_runtime2.jsx)(InputGroupText_default, {
  children: (0, import_jsx_runtime2.jsx)(FormCheckInput_default, {
    type: "radio",
    ...props
  })
});
var InputGroup = React2.forwardRef(({
  bsPrefix,
  size,
  hasValidation,
  className,
  // Need to define the default "as" during prop destructuring to be compatible with styled-components github.com/react-bootstrap/react-bootstrap/issues/3595
  as: Component = "div",
  ...props
}, ref) => {
  bsPrefix = useBootstrapPrefix(bsPrefix, "input-group");
  const contextValue = (0, import_react.useMemo)(() => ({}), []);
  return (0, import_jsx_runtime2.jsx)(InputGroupContext_default.Provider, {
    value: contextValue,
    children: (0, import_jsx_runtime2.jsx)(Component, {
      ref,
      ...props,
      className: (0, import_classnames2.default)(className, bsPrefix, size && `${bsPrefix}-${size}`, hasValidation && "has-validation")
    })
  });
});
InputGroup.displayName = "InputGroup";
var InputGroup_default = Object.assign(InputGroup, {
  Text: InputGroupText_default,
  Radio: InputGroupRadio,
  Checkbox: InputGroupCheckbox
});
export {
  InputGroup_default as default
};
//# sourceMappingURL=react-bootstrap_InputGroup.js.map
