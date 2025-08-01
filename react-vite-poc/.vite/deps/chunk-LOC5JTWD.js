import {
  require_classnames,
  require_jsx_runtime,
  useBootstrapPrefix
} from "./chunk-5TWWBDIN.js";
import {
  __toESM,
  require_react
} from "./chunk-5WQJO2FO.js";

// node_modules/react-bootstrap/esm/FormCheckInput.js
var import_classnames = __toESM(require_classnames());
var React2 = __toESM(require_react());
var import_react = __toESM(require_react());

// node_modules/react-bootstrap/esm/FormContext.js
var React = __toESM(require_react());
var FormContext = React.createContext({});
var FormContext_default = FormContext;

// node_modules/react-bootstrap/esm/FormCheckInput.js
var import_jsx_runtime = __toESM(require_jsx_runtime());
var FormCheckInput = React2.forwardRef(({
  id,
  bsPrefix,
  className,
  type = "checkbox",
  isValid = false,
  isInvalid = false,
  // Need to define the default "as" during prop destructuring to be compatible with styled-components github.com/react-bootstrap/react-bootstrap/issues/3595
  as: Component = "input",
  ...props
}, ref) => {
  const {
    controlId
  } = (0, import_react.useContext)(FormContext_default);
  bsPrefix = useBootstrapPrefix(bsPrefix, "form-check-input");
  return (0, import_jsx_runtime.jsx)(Component, {
    ...props,
    ref,
    type,
    id: id || controlId,
    className: (0, import_classnames.default)(className, bsPrefix, isValid && "is-valid", isInvalid && "is-invalid")
  });
});
FormCheckInput.displayName = "FormCheckInput";
var FormCheckInput_default = FormCheckInput;

export {
  FormContext_default,
  FormCheckInput_default
};
//# sourceMappingURL=chunk-LOC5JTWD.js.map
