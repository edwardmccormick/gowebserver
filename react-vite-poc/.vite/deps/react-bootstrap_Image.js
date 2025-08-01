"use client";
import {
  require_prop_types
} from "./chunk-L3ZLE5PF.js";
import {
  require_classnames,
  require_jsx_runtime,
  useBootstrapPrefix
} from "./chunk-5TWWBDIN.js";
import {
  __toESM,
  require_react
} from "./chunk-5WQJO2FO.js";

// node_modules/react-bootstrap/esm/Image.js
var import_classnames = __toESM(require_classnames());
var React = __toESM(require_react());
var import_prop_types = __toESM(require_prop_types());
var import_jsx_runtime = __toESM(require_jsx_runtime());
var propTypes = {
  /**
   * @default 'img'
   */
  bsPrefix: import_prop_types.default.string,
  /**
   * Sets image as fluid image.
   */
  fluid: import_prop_types.default.bool,
  /**
   * Sets image shape as rounded.
   */
  rounded: import_prop_types.default.bool,
  /**
   * Sets image shape as circle.
   */
  roundedCircle: import_prop_types.default.bool,
  /**
   * Sets image shape as thumbnail.
   */
  thumbnail: import_prop_types.default.bool
};
var Image = React.forwardRef(({
  bsPrefix,
  className,
  fluid = false,
  rounded = false,
  roundedCircle = false,
  thumbnail = false,
  ...props
}, ref) => {
  bsPrefix = useBootstrapPrefix(bsPrefix, "img");
  return (0, import_jsx_runtime.jsx)("img", {
    ref,
    ...props,
    className: (0, import_classnames.default)(className, fluid && `${bsPrefix}-fluid`, rounded && `rounded`, roundedCircle && `rounded-circle`, thumbnail && `${bsPrefix}-thumbnail`)
  });
});
Image.displayName = "Image";
var Image_default = Image;
export {
  Image_default as default,
  propTypes
};
//# sourceMappingURL=react-bootstrap_Image.js.map
