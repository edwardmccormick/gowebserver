import {
  require_warning
} from "./chunk-76RW7NSA.js";
import {
  hasChildOfType
} from "./chunk-2A7GM3I4.js";
import {
  FormCheckInput_default,
  FormContext_default
} from "./chunk-LOC5JTWD.js";
import {
  require_prop_types
} from "./chunk-L3ZLE5PF.js";
import {
  require_classnames,
  require_jsx_runtime,
  useBootstrapBreakpoints,
  useBootstrapMinBreakpoint,
  useBootstrapPrefix
} from "./chunk-5TWWBDIN.js";
import {
  __toESM,
  require_react
} from "./chunk-5WQJO2FO.js";

// node_modules/react-bootstrap/esm/Form.js
var import_classnames12 = __toESM(require_classnames());
var import_prop_types2 = __toESM(require_prop_types());
var React14 = __toESM(require_react());

// node_modules/react-bootstrap/esm/FormCheck.js
var import_classnames3 = __toESM(require_classnames());
var React3 = __toESM(require_react());
var import_react2 = __toESM(require_react());

// node_modules/react-bootstrap/esm/Feedback.js
var import_classnames = __toESM(require_classnames());
var React = __toESM(require_react());
var import_prop_types = __toESM(require_prop_types());
var import_jsx_runtime = __toESM(require_jsx_runtime());
var propTypes = {
  /**
   * Specify whether the feedback is for valid or invalid fields
   *
   * @type {('valid'|'invalid')}
   */
  type: import_prop_types.default.string,
  /** Display feedback as a tooltip. */
  tooltip: import_prop_types.default.bool,
  as: import_prop_types.default.elementType
};
var Feedback = React.forwardRef(
  // Need to define the default "as" during prop destructuring to be compatible with styled-components github.com/react-bootstrap/react-bootstrap/issues/3595
  ({
    as: Component = "div",
    className,
    type = "valid",
    tooltip = false,
    ...props
  }, ref) => (0, import_jsx_runtime.jsx)(Component, {
    ...props,
    ref,
    className: (0, import_classnames.default)(className, `${type}-${tooltip ? "tooltip" : "feedback"}`)
  })
);
Feedback.displayName = "Feedback";
Feedback.propTypes = propTypes;
var Feedback_default = Feedback;

// node_modules/react-bootstrap/esm/FormCheckLabel.js
var import_classnames2 = __toESM(require_classnames());
var React2 = __toESM(require_react());
var import_react = __toESM(require_react());
var import_jsx_runtime2 = __toESM(require_jsx_runtime());
var FormCheckLabel = React2.forwardRef(({
  bsPrefix,
  className,
  htmlFor,
  ...props
}, ref) => {
  const {
    controlId
  } = (0, import_react.useContext)(FormContext_default);
  bsPrefix = useBootstrapPrefix(bsPrefix, "form-check-label");
  return (0, import_jsx_runtime2.jsx)("label", {
    ...props,
    ref,
    htmlFor: htmlFor || controlId,
    className: (0, import_classnames2.default)(className, bsPrefix)
  });
});
FormCheckLabel.displayName = "FormCheckLabel";
var FormCheckLabel_default = FormCheckLabel;

// node_modules/react-bootstrap/esm/FormCheck.js
var import_jsx_runtime3 = __toESM(require_jsx_runtime());
var import_jsx_runtime4 = __toESM(require_jsx_runtime());
var import_jsx_runtime5 = __toESM(require_jsx_runtime());
var FormCheck = React3.forwardRef(({
  id,
  bsPrefix,
  bsSwitchPrefix,
  inline = false,
  reverse = false,
  disabled = false,
  isValid = false,
  isInvalid = false,
  feedbackTooltip = false,
  feedback,
  feedbackType,
  className,
  style,
  title = "",
  type = "checkbox",
  label,
  children,
  // Need to define the default "as" during prop destructuring to be compatible with styled-components github.com/react-bootstrap/react-bootstrap/issues/3595
  as = "input",
  ...props
}, ref) => {
  bsPrefix = useBootstrapPrefix(bsPrefix, "form-check");
  bsSwitchPrefix = useBootstrapPrefix(bsSwitchPrefix, "form-switch");
  const {
    controlId
  } = (0, import_react2.useContext)(FormContext_default);
  const innerFormContext = (0, import_react2.useMemo)(() => ({
    controlId: id || controlId
  }), [controlId, id]);
  const hasLabel = !children && label != null && label !== false || hasChildOfType(children, FormCheckLabel_default);
  const input = (0, import_jsx_runtime3.jsx)(FormCheckInput_default, {
    ...props,
    type: type === "switch" ? "checkbox" : type,
    ref,
    isValid,
    isInvalid,
    disabled,
    as
  });
  return (0, import_jsx_runtime3.jsx)(FormContext_default.Provider, {
    value: innerFormContext,
    children: (0, import_jsx_runtime3.jsx)("div", {
      style,
      className: (0, import_classnames3.default)(className, hasLabel && bsPrefix, inline && `${bsPrefix}-inline`, reverse && `${bsPrefix}-reverse`, type === "switch" && bsSwitchPrefix),
      children: children || (0, import_jsx_runtime5.jsxs)(import_jsx_runtime4.Fragment, {
        children: [input, hasLabel && (0, import_jsx_runtime3.jsx)(FormCheckLabel_default, {
          title,
          children: label
        }), feedback && (0, import_jsx_runtime3.jsx)(Feedback_default, {
          type: feedbackType,
          tooltip: feedbackTooltip,
          children: feedback
        })]
      })
    })
  });
});
FormCheck.displayName = "FormCheck";
var FormCheck_default = Object.assign(FormCheck, {
  Input: FormCheckInput_default,
  Label: FormCheckLabel_default
});

// node_modules/react-bootstrap/esm/FormControl.js
var import_classnames4 = __toESM(require_classnames());
var React4 = __toESM(require_react());
var import_react3 = __toESM(require_react());
var import_warning = __toESM(require_warning());
var import_jsx_runtime6 = __toESM(require_jsx_runtime());
var FormControl = React4.forwardRef(({
  bsPrefix,
  type,
  size,
  htmlSize,
  id,
  className,
  isValid = false,
  isInvalid = false,
  plaintext,
  readOnly,
  // Need to define the default "as" during prop destructuring to be compatible with styled-components github.com/react-bootstrap/react-bootstrap/issues/3595
  as: Component = "input",
  ...props
}, ref) => {
  const {
    controlId
  } = (0, import_react3.useContext)(FormContext_default);
  bsPrefix = useBootstrapPrefix(bsPrefix, "form-control");
  true ? (0, import_warning.default)(controlId == null || !id, "`controlId` is ignored on `<FormControl>` when `id` is specified.") : void 0;
  return (0, import_jsx_runtime6.jsx)(Component, {
    ...props,
    type,
    size: htmlSize,
    ref,
    readOnly,
    id: id || controlId,
    className: (0, import_classnames4.default)(className, plaintext ? `${bsPrefix}-plaintext` : bsPrefix, size && `${bsPrefix}-${size}`, type === "color" && `${bsPrefix}-color`, isValid && "is-valid", isInvalid && "is-invalid")
  });
});
FormControl.displayName = "FormControl";
var FormControl_default = Object.assign(FormControl, {
  Feedback: Feedback_default
});

// node_modules/react-bootstrap/esm/FormFloating.js
var React5 = __toESM(require_react());
var import_classnames5 = __toESM(require_classnames());
var import_jsx_runtime7 = __toESM(require_jsx_runtime());
var FormFloating = React5.forwardRef(({
  className,
  bsPrefix,
  as: Component = "div",
  ...props
}, ref) => {
  bsPrefix = useBootstrapPrefix(bsPrefix, "form-floating");
  return (0, import_jsx_runtime7.jsx)(Component, {
    ref,
    className: (0, import_classnames5.default)(className, bsPrefix),
    ...props
  });
});
FormFloating.displayName = "FormFloating";
var FormFloating_default = FormFloating;

// node_modules/react-bootstrap/esm/FormGroup.js
var React6 = __toESM(require_react());
var import_react4 = __toESM(require_react());
var import_jsx_runtime8 = __toESM(require_jsx_runtime());
var FormGroup = React6.forwardRef(({
  controlId,
  // Need to define the default "as" during prop destructuring to be compatible with styled-components github.com/react-bootstrap/react-bootstrap/issues/3595
  as: Component = "div",
  ...props
}, ref) => {
  const context = (0, import_react4.useMemo)(() => ({
    controlId
  }), [controlId]);
  return (0, import_jsx_runtime8.jsx)(FormContext_default.Provider, {
    value: context,
    children: (0, import_jsx_runtime8.jsx)(Component, {
      ...props,
      ref
    })
  });
});
FormGroup.displayName = "FormGroup";
var FormGroup_default = FormGroup;

// node_modules/react-bootstrap/esm/FormLabel.js
var import_classnames7 = __toESM(require_classnames());
var React8 = __toESM(require_react());
var import_react5 = __toESM(require_react());
var import_warning2 = __toESM(require_warning());

// node_modules/react-bootstrap/esm/Col.js
var import_classnames6 = __toESM(require_classnames());
var React7 = __toESM(require_react());
var import_jsx_runtime9 = __toESM(require_jsx_runtime());
function useCol({
  as,
  bsPrefix,
  className,
  ...props
}) {
  bsPrefix = useBootstrapPrefix(bsPrefix, "col");
  const breakpoints = useBootstrapBreakpoints();
  const minBreakpoint = useBootstrapMinBreakpoint();
  const spans = [];
  const classes = [];
  breakpoints.forEach((brkPoint) => {
    const propValue = props[brkPoint];
    delete props[brkPoint];
    let span;
    let offset;
    let order;
    if (typeof propValue === "object" && propValue != null) {
      ({
        span,
        offset,
        order
      } = propValue);
    } else {
      span = propValue;
    }
    const infix = brkPoint !== minBreakpoint ? `-${brkPoint}` : "";
    if (span) spans.push(span === true ? `${bsPrefix}${infix}` : `${bsPrefix}${infix}-${span}`);
    if (order != null) classes.push(`order${infix}-${order}`);
    if (offset != null) classes.push(`offset${infix}-${offset}`);
  });
  return [{
    ...props,
    className: (0, import_classnames6.default)(className, ...spans, ...classes)
  }, {
    as,
    bsPrefix,
    spans
  }];
}
var Col = React7.forwardRef(
  // Need to define the default "as" during prop destructuring to be compatible with styled-components github.com/react-bootstrap/react-bootstrap/issues/3595
  (props, ref) => {
    const [{
      className,
      ...colProps
    }, {
      as: Component = "div",
      bsPrefix,
      spans
    }] = useCol(props);
    return (0, import_jsx_runtime9.jsx)(Component, {
      ...colProps,
      ref,
      className: (0, import_classnames6.default)(className, !spans.length && bsPrefix)
    });
  }
);
Col.displayName = "Col";
var Col_default = Col;

// node_modules/react-bootstrap/esm/FormLabel.js
var import_jsx_runtime10 = __toESM(require_jsx_runtime());
var FormLabel = React8.forwardRef(({
  // Need to define the default "as" during prop destructuring to be compatible with styled-components github.com/react-bootstrap/react-bootstrap/issues/3595
  as: Component = "label",
  bsPrefix,
  column = false,
  visuallyHidden = false,
  className,
  htmlFor,
  ...props
}, ref) => {
  const {
    controlId
  } = (0, import_react5.useContext)(FormContext_default);
  bsPrefix = useBootstrapPrefix(bsPrefix, "form-label");
  let columnClass = "col-form-label";
  if (typeof column === "string") columnClass = `${columnClass} ${columnClass}-${column}`;
  const classes = (0, import_classnames7.default)(className, bsPrefix, visuallyHidden && "visually-hidden", column && columnClass);
  true ? (0, import_warning2.default)(controlId == null || !htmlFor, "`controlId` is ignored on `<FormLabel>` when `htmlFor` is specified.") : void 0;
  htmlFor = htmlFor || controlId;
  if (column) return (0, import_jsx_runtime10.jsx)(Col_default, {
    ref,
    as: "label",
    className: classes,
    htmlFor,
    ...props
  });
  return (0, import_jsx_runtime10.jsx)(Component, {
    ref,
    className: classes,
    htmlFor,
    ...props
  });
});
FormLabel.displayName = "FormLabel";
var FormLabel_default = FormLabel;

// node_modules/react-bootstrap/esm/FormRange.js
var import_classnames8 = __toESM(require_classnames());
var React9 = __toESM(require_react());
var import_react6 = __toESM(require_react());
var import_jsx_runtime11 = __toESM(require_jsx_runtime());
var FormRange = React9.forwardRef(({
  bsPrefix,
  className,
  id,
  ...props
}, ref) => {
  const {
    controlId
  } = (0, import_react6.useContext)(FormContext_default);
  bsPrefix = useBootstrapPrefix(bsPrefix, "form-range");
  return (0, import_jsx_runtime11.jsx)("input", {
    ...props,
    type: "range",
    ref,
    className: (0, import_classnames8.default)(className, bsPrefix),
    id: id || controlId
  });
});
FormRange.displayName = "FormRange";
var FormRange_default = FormRange;

// node_modules/react-bootstrap/esm/FormSelect.js
var import_classnames9 = __toESM(require_classnames());
var React10 = __toESM(require_react());
var import_react7 = __toESM(require_react());
var import_jsx_runtime12 = __toESM(require_jsx_runtime());
var FormSelect = React10.forwardRef(({
  bsPrefix,
  size,
  htmlSize,
  className,
  isValid = false,
  isInvalid = false,
  id,
  ...props
}, ref) => {
  const {
    controlId
  } = (0, import_react7.useContext)(FormContext_default);
  bsPrefix = useBootstrapPrefix(bsPrefix, "form-select");
  return (0, import_jsx_runtime12.jsx)("select", {
    ...props,
    size: htmlSize,
    ref,
    className: (0, import_classnames9.default)(className, bsPrefix, size && `${bsPrefix}-${size}`, isValid && `is-valid`, isInvalid && `is-invalid`),
    id: id || controlId
  });
});
FormSelect.displayName = "FormSelect";
var FormSelect_default = FormSelect;

// node_modules/react-bootstrap/esm/FormText.js
var import_classnames10 = __toESM(require_classnames());
var React11 = __toESM(require_react());
var import_jsx_runtime13 = __toESM(require_jsx_runtime());
var FormText = React11.forwardRef(
  // Need to define the default "as" during prop destructuring to be compatible with styled-components github.com/react-bootstrap/react-bootstrap/issues/3595
  ({
    bsPrefix,
    className,
    as: Component = "small",
    muted,
    ...props
  }, ref) => {
    bsPrefix = useBootstrapPrefix(bsPrefix, "form-text");
    return (0, import_jsx_runtime13.jsx)(Component, {
      ...props,
      ref,
      className: (0, import_classnames10.default)(className, bsPrefix, muted && "text-muted")
    });
  }
);
FormText.displayName = "FormText";
var FormText_default = FormText;

// node_modules/react-bootstrap/esm/Switch.js
var React12 = __toESM(require_react());
var import_jsx_runtime14 = __toESM(require_jsx_runtime());
var Switch = React12.forwardRef((props, ref) => (0, import_jsx_runtime14.jsx)(FormCheck_default, {
  ...props,
  ref,
  type: "switch"
}));
Switch.displayName = "Switch";
var Switch_default = Object.assign(Switch, {
  Input: FormCheck_default.Input,
  Label: FormCheck_default.Label
});

// node_modules/react-bootstrap/esm/FloatingLabel.js
var import_classnames11 = __toESM(require_classnames());
var React13 = __toESM(require_react());
var import_jsx_runtime15 = __toESM(require_jsx_runtime());
var import_jsx_runtime16 = __toESM(require_jsx_runtime());
var FloatingLabel = React13.forwardRef(({
  bsPrefix,
  className,
  children,
  controlId,
  label,
  ...props
}, ref) => {
  bsPrefix = useBootstrapPrefix(bsPrefix, "form-floating");
  return (0, import_jsx_runtime16.jsxs)(FormGroup_default, {
    ref,
    className: (0, import_classnames11.default)(className, bsPrefix),
    controlId,
    ...props,
    children: [children, (0, import_jsx_runtime15.jsx)("label", {
      htmlFor: controlId,
      children: label
    })]
  });
});
FloatingLabel.displayName = "FloatingLabel";
var FloatingLabel_default = FloatingLabel;

// node_modules/react-bootstrap/esm/Form.js
var import_jsx_runtime17 = __toESM(require_jsx_runtime());
var propTypes2 = {
  /**
   * The Form `ref` will be forwarded to the underlying element,
   * which means, unless it's rendered `as` a composite component,
   * it will be a DOM node, when resolved.
   *
   * @type {ReactRef}
   * @alias ref
   */
  _ref: import_prop_types2.default.any,
  /**
   * Mark a form as having been validated. Setting it to `true` will
   * toggle any validation styles on the forms elements.
   */
  validated: import_prop_types2.default.bool,
  as: import_prop_types2.default.elementType
};
var Form = React14.forwardRef(({
  className,
  validated,
  // Need to define the default "as" during prop destructuring to be compatible with styled-components github.com/react-bootstrap/react-bootstrap/issues/3595
  as: Component = "form",
  ...props
}, ref) => (0, import_jsx_runtime17.jsx)(Component, {
  ...props,
  ref,
  className: (0, import_classnames12.default)(className, validated && "was-validated")
}));
Form.displayName = "Form";
Form.propTypes = propTypes2;
var Form_default = Object.assign(Form, {
  Group: FormGroup_default,
  Control: FormControl_default,
  Floating: FormFloating_default,
  Check: FormCheck_default,
  Switch: Switch_default,
  Label: FormLabel_default,
  Text: FormText_default,
  Range: FormRange_default,
  Select: FormSelect_default,
  FloatingLabel: FloatingLabel_default
});
export {
  Form_default as default
};
//# sourceMappingURL=react-bootstrap_Form.js.map
