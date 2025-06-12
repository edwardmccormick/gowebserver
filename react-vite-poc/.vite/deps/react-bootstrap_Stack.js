"use client";
import {
  require_prop_types
} from "./chunk-L3ZLE5PF.js";
import {
  DEFAULT_BREAKPOINTS,
  DEFAULT_MIN_BREAKPOINT,
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

// node_modules/react-bootstrap/esm/Stack.js
var import_classnames = __toESM(require_classnames());
var React = __toESM(require_react());

// node_modules/react-bootstrap/esm/createUtilityClasses.js
var import_prop_types = __toESM(require_prop_types());
function createUtilityClassName(utilityValues, breakpoints = DEFAULT_BREAKPOINTS, minBreakpoint = DEFAULT_MIN_BREAKPOINT) {
  const classes = [];
  Object.entries(utilityValues).forEach(([utilName, utilValue]) => {
    if (utilValue != null) {
      if (typeof utilValue === "object") {
        breakpoints.forEach((brkPoint) => {
          const bpValue = utilValue[brkPoint];
          if (bpValue != null) {
            const infix = brkPoint !== minBreakpoint ? `-${brkPoint}` : "";
            classes.push(`${utilName}${infix}-${bpValue}`);
          }
        });
      } else {
        classes.push(`${utilName}-${utilValue}`);
      }
    }
  });
  return classes;
}

// node_modules/react-bootstrap/esm/Stack.js
var import_jsx_runtime = __toESM(require_jsx_runtime());
var Stack = React.forwardRef(({
  as: Component = "div",
  bsPrefix,
  className,
  direction,
  gap,
  ...props
}, ref) => {
  bsPrefix = useBootstrapPrefix(bsPrefix, direction === "horizontal" ? "hstack" : "vstack");
  const breakpoints = useBootstrapBreakpoints();
  const minBreakpoint = useBootstrapMinBreakpoint();
  return (0, import_jsx_runtime.jsx)(Component, {
    ...props,
    ref,
    className: (0, import_classnames.default)(className, bsPrefix, ...createUtilityClassName({
      gap
    }, breakpoints, minBreakpoint))
  });
});
Stack.displayName = "Stack";
var Stack_default = Stack;
export {
  Stack_default as default
};
//# sourceMappingURL=react-bootstrap_Stack.js.map
