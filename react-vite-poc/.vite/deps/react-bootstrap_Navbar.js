"use client";
import {
  useIsomorphicEffect_default
} from "./chunk-ZERQRB3V.js";
import {
  AbstractModalHeader_default,
  BootstrapModalManager_default,
  Fade_default,
  ModalContext_default,
  Modal_default,
  divWithClassName_default,
  getSharedManager
} from "./chunk-6BEFPC2L.js";
import "./chunk-SQ4NDZ5V.js";
import {
  useEventCallback
} from "./chunk-PJQH5SKO.js";
import {
  Collapse_default
} from "./chunk-D6ZGIBSS.js";
import {
  getChildRef
} from "./chunk-MTDRLTFK.js";
import {
  ENTERED,
  ENTERING,
  EXITING,
  TransitionWrapper_default,
  transitionEndListener
} from "./chunk-T2ES4POL.js";
import "./chunk-S2TLU4L2.js";
import "./chunk-QOVZGY7A.js";
import "./chunk-L3ZLE5PF.js";
import {
  NavbarContext_default
} from "./chunk-3BORRP6Y.js";
import {
  SelectableContext_default
} from "./chunk-25OZPJTB.js";
import "./chunk-KDSTLSK6.js";
import "./chunk-W3GZWJ6O.js";
import "./chunk-TZNIXWST.js";
import "./chunk-ZEXTMIKL.js";
import {
  useUncontrolled
} from "./chunk-BGOLUZHQ.js";
import "./chunk-EEJ6RG5R.js";
import {
  require_classnames,
  require_jsx_runtime,
  useBootstrapPrefix
} from "./chunk-5TWWBDIN.js";
import {
  __toESM,
  require_react
} from "./chunk-5WQJO2FO.js";

// node_modules/react-bootstrap/esm/Navbar.js
var import_classnames9 = __toESM(require_classnames());
var React11 = __toESM(require_react());
var import_react7 = __toESM(require_react());

// node_modules/react-bootstrap/esm/NavbarBrand.js
var import_classnames = __toESM(require_classnames());
var React = __toESM(require_react());
var import_jsx_runtime = __toESM(require_jsx_runtime());
var NavbarBrand = React.forwardRef(({
  bsPrefix,
  className,
  as,
  ...props
}, ref) => {
  bsPrefix = useBootstrapPrefix(bsPrefix, "navbar-brand");
  const Component = as || (props.href ? "a" : "span");
  return (0, import_jsx_runtime.jsx)(Component, {
    ...props,
    ref,
    className: (0, import_classnames.default)(className, bsPrefix)
  });
});
NavbarBrand.displayName = "NavbarBrand";
var NavbarBrand_default = NavbarBrand;

// node_modules/react-bootstrap/esm/NavbarCollapse.js
var React2 = __toESM(require_react());
var import_react = __toESM(require_react());
var import_jsx_runtime2 = __toESM(require_jsx_runtime());
var NavbarCollapse = React2.forwardRef(({
  children,
  bsPrefix,
  ...props
}, ref) => {
  bsPrefix = useBootstrapPrefix(bsPrefix, "navbar-collapse");
  const context = (0, import_react.useContext)(NavbarContext_default);
  return (0, import_jsx_runtime2.jsx)(Collapse_default, {
    in: !!(context && context.expanded),
    ...props,
    children: (0, import_jsx_runtime2.jsx)("div", {
      ref,
      className: bsPrefix,
      children
    })
  });
});
NavbarCollapse.displayName = "NavbarCollapse";
var NavbarCollapse_default = NavbarCollapse;

// node_modules/react-bootstrap/esm/NavbarToggle.js
var import_classnames2 = __toESM(require_classnames());
var React3 = __toESM(require_react());
var import_react2 = __toESM(require_react());
var import_jsx_runtime3 = __toESM(require_jsx_runtime());
var NavbarToggle = React3.forwardRef(({
  bsPrefix,
  className,
  children,
  label = "Toggle navigation",
  // Need to define the default "as" during prop destructuring to be compatible with styled-components github.com/react-bootstrap/react-bootstrap/issues/3595
  as: Component = "button",
  onClick,
  ...props
}, ref) => {
  bsPrefix = useBootstrapPrefix(bsPrefix, "navbar-toggler");
  const {
    onToggle,
    expanded
  } = (0, import_react2.useContext)(NavbarContext_default) || {};
  const handleClick = useEventCallback((e) => {
    if (onClick) onClick(e);
    if (onToggle) onToggle();
  });
  if (Component === "button") {
    props.type = "button";
  }
  return (0, import_jsx_runtime3.jsx)(Component, {
    ...props,
    ref,
    onClick: handleClick,
    "aria-label": label,
    className: (0, import_classnames2.default)(className, bsPrefix, !expanded && "collapsed"),
    children: children || (0, import_jsx_runtime3.jsx)("span", {
      className: `${bsPrefix}-icon`
    })
  });
});
NavbarToggle.displayName = "NavbarToggle";
var NavbarToggle_default = NavbarToggle;

// node_modules/react-bootstrap/esm/NavbarOffcanvas.js
var React9 = __toESM(require_react());
var import_react6 = __toESM(require_react());

// node_modules/react-bootstrap/esm/Offcanvas.js
var import_classnames7 = __toESM(require_classnames());

// node_modules/@restart/hooks/esm/useMediaQuery.js
var import_react3 = __toESM(require_react());
var matchersByWindow = /* @__PURE__ */ new WeakMap();
var getMatcher = (query, targetWindow) => {
  if (!query || !targetWindow) return void 0;
  const matchers = matchersByWindow.get(targetWindow) || /* @__PURE__ */ new Map();
  matchersByWindow.set(targetWindow, matchers);
  let mql = matchers.get(query);
  if (!mql) {
    mql = targetWindow.matchMedia(query);
    mql.refCount = 0;
    matchers.set(mql.media, mql);
  }
  return mql;
};
function useMediaQuery(query, targetWindow = typeof window === "undefined" ? void 0 : window) {
  const mql = getMatcher(query, targetWindow);
  const [matches, setMatches] = (0, import_react3.useState)(() => mql ? mql.matches : false);
  useIsomorphicEffect_default(() => {
    let mql2 = getMatcher(query, targetWindow);
    if (!mql2) {
      return setMatches(false);
    }
    let matchers = matchersByWindow.get(targetWindow);
    const handleChange = () => {
      setMatches(mql2.matches);
    };
    mql2.refCount++;
    mql2.addListener(handleChange);
    handleChange();
    return () => {
      mql2.removeListener(handleChange);
      mql2.refCount--;
      if (mql2.refCount <= 0) {
        matchers == null ? void 0 : matchers.delete(mql2.media);
      }
      mql2 = void 0;
    };
  }, [query]);
  return matches;
}

// node_modules/@restart/hooks/esm/useBreakpoint.js
var import_react4 = __toESM(require_react());
function createBreakpointHook(breakpointValues) {
  const names = Object.keys(breakpointValues);
  function and(query, next) {
    if (query === next) {
      return next;
    }
    return query ? `${query} and ${next}` : next;
  }
  function getNext(breakpoint) {
    return names[Math.min(names.indexOf(breakpoint) + 1, names.length - 1)];
  }
  function getMaxQuery(breakpoint) {
    const next = getNext(breakpoint);
    let value = breakpointValues[next];
    if (typeof value === "number") value = `${value - 0.2}px`;
    else value = `calc(${value} - 0.2px)`;
    return `(max-width: ${value})`;
  }
  function getMinQuery(breakpoint) {
    let value = breakpointValues[breakpoint];
    if (typeof value === "number") {
      value = `${value}px`;
    }
    return `(min-width: ${value})`;
  }
  function useBreakpoint2(breakpointOrMap, direction, window2) {
    let breakpointMap;
    if (typeof breakpointOrMap === "object") {
      breakpointMap = breakpointOrMap;
      window2 = direction;
      direction = true;
    } else {
      direction = direction || true;
      breakpointMap = {
        [breakpointOrMap]: direction
      };
    }
    let query = (0, import_react4.useMemo)(() => Object.entries(breakpointMap).reduce((query2, [key, direction2]) => {
      if (direction2 === "up" || direction2 === true) {
        query2 = and(query2, getMinQuery(key));
      }
      if (direction2 === "down" || direction2 === true) {
        query2 = and(query2, getMaxQuery(key));
      }
      return query2;
    }, ""), [JSON.stringify(breakpointMap)]);
    return useMediaQuery(query, window2);
  }
  return useBreakpoint2;
}
var useBreakpoint = createBreakpointHook({
  xs: 0,
  sm: 576,
  md: 768,
  lg: 992,
  xl: 1200,
  xxl: 1400
});
var useBreakpoint_default = useBreakpoint;

// node_modules/react-bootstrap/esm/Offcanvas.js
var React8 = __toESM(require_react());
var import_react5 = __toESM(require_react());

// node_modules/react-bootstrap/esm/OffcanvasBody.js
var React4 = __toESM(require_react());
var import_classnames3 = __toESM(require_classnames());
var import_jsx_runtime4 = __toESM(require_jsx_runtime());
var OffcanvasBody = React4.forwardRef(({
  className,
  bsPrefix,
  as: Component = "div",
  ...props
}, ref) => {
  bsPrefix = useBootstrapPrefix(bsPrefix, "offcanvas-body");
  return (0, import_jsx_runtime4.jsx)(Component, {
    ref,
    className: (0, import_classnames3.default)(className, bsPrefix),
    ...props
  });
});
OffcanvasBody.displayName = "OffcanvasBody";
var OffcanvasBody_default = OffcanvasBody;

// node_modules/react-bootstrap/esm/OffcanvasToggling.js
var import_classnames4 = __toESM(require_classnames());
var React5 = __toESM(require_react());
var import_jsx_runtime5 = __toESM(require_jsx_runtime());
var transitionStyles = {
  [ENTERING]: "show",
  [ENTERED]: "show"
};
var OffcanvasToggling = React5.forwardRef(({
  bsPrefix,
  className,
  children,
  in: inProp = false,
  mountOnEnter = false,
  unmountOnExit = false,
  appear = false,
  ...props
}, ref) => {
  bsPrefix = useBootstrapPrefix(bsPrefix, "offcanvas");
  return (0, import_jsx_runtime5.jsx)(TransitionWrapper_default, {
    ref,
    addEndListener: transitionEndListener,
    in: inProp,
    mountOnEnter,
    unmountOnExit,
    appear,
    ...props,
    childRef: getChildRef(children),
    children: (status, innerProps) => React5.cloneElement(children, {
      ...innerProps,
      className: (0, import_classnames4.default)(className, children.props.className, (status === ENTERING || status === EXITING) && `${bsPrefix}-toggling`, transitionStyles[status])
    })
  });
});
OffcanvasToggling.displayName = "OffcanvasToggling";
var OffcanvasToggling_default = OffcanvasToggling;

// node_modules/react-bootstrap/esm/OffcanvasHeader.js
var import_classnames5 = __toESM(require_classnames());
var React6 = __toESM(require_react());
var import_jsx_runtime6 = __toESM(require_jsx_runtime());
var OffcanvasHeader = React6.forwardRef(({
  bsPrefix,
  className,
  closeLabel = "Close",
  closeButton = false,
  ...props
}, ref) => {
  bsPrefix = useBootstrapPrefix(bsPrefix, "offcanvas-header");
  return (0, import_jsx_runtime6.jsx)(AbstractModalHeader_default, {
    ref,
    ...props,
    className: (0, import_classnames5.default)(className, bsPrefix),
    closeLabel,
    closeButton
  });
});
OffcanvasHeader.displayName = "OffcanvasHeader";
var OffcanvasHeader_default = OffcanvasHeader;

// node_modules/react-bootstrap/esm/OffcanvasTitle.js
var React7 = __toESM(require_react());
var import_classnames6 = __toESM(require_classnames());
var import_jsx_runtime7 = __toESM(require_jsx_runtime());
var DivStyledAsH5 = divWithClassName_default("h5");
var OffcanvasTitle = React7.forwardRef(({
  className,
  bsPrefix,
  as: Component = DivStyledAsH5,
  ...props
}, ref) => {
  bsPrefix = useBootstrapPrefix(bsPrefix, "offcanvas-title");
  return (0, import_jsx_runtime7.jsx)(Component, {
    ref,
    className: (0, import_classnames6.default)(className, bsPrefix),
    ...props
  });
});
OffcanvasTitle.displayName = "OffcanvasTitle";
var OffcanvasTitle_default = OffcanvasTitle;

// node_modules/react-bootstrap/esm/Offcanvas.js
var import_jsx_runtime8 = __toESM(require_jsx_runtime());
var import_jsx_runtime9 = __toESM(require_jsx_runtime());
var import_jsx_runtime10 = __toESM(require_jsx_runtime());
function DialogTransition(props) {
  return (0, import_jsx_runtime8.jsx)(OffcanvasToggling_default, {
    ...props
  });
}
function BackdropTransition(props) {
  return (0, import_jsx_runtime8.jsx)(Fade_default, {
    ...props
  });
}
var Offcanvas = React8.forwardRef(({
  bsPrefix,
  className,
  children,
  "aria-labelledby": ariaLabelledby,
  placement = "start",
  responsive,
  /* BaseModal props */
  show = false,
  backdrop = true,
  keyboard = true,
  scroll = false,
  onEscapeKeyDown,
  onShow,
  onHide,
  container,
  autoFocus = true,
  enforceFocus = true,
  restoreFocus = true,
  restoreFocusOptions,
  onEntered,
  onExit,
  onExiting,
  onEnter,
  onEntering,
  onExited,
  backdropClassName,
  manager: propsManager,
  renderStaticNode = false,
  ...props
}, ref) => {
  const modalManager = (0, import_react5.useRef)();
  bsPrefix = useBootstrapPrefix(bsPrefix, "offcanvas");
  const [showOffcanvas, setShowOffcanvas] = (0, import_react5.useState)(false);
  const handleHide = useEventCallback(onHide);
  const hideResponsiveOffcanvas = useBreakpoint_default(responsive || "xs", "up");
  (0, import_react5.useEffect)(() => {
    setShowOffcanvas(responsive ? show && !hideResponsiveOffcanvas : show);
  }, [show, responsive, hideResponsiveOffcanvas]);
  const modalContext = (0, import_react5.useMemo)(() => ({
    onHide: handleHide
  }), [handleHide]);
  function getModalManager() {
    if (propsManager) return propsManager;
    if (scroll) {
      if (!modalManager.current) modalManager.current = new BootstrapModalManager_default({
        handleContainerOverflow: false
      });
      return modalManager.current;
    }
    return getSharedManager();
  }
  const handleEnter = (node, ...args) => {
    if (node) node.style.visibility = "visible";
    onEnter == null || onEnter(node, ...args);
  };
  const handleExited = (node, ...args) => {
    if (node) node.style.visibility = "";
    onExited == null || onExited(...args);
  };
  const renderBackdrop = (0, import_react5.useCallback)((backdropProps) => (0, import_jsx_runtime8.jsx)("div", {
    ...backdropProps,
    className: (0, import_classnames7.default)(`${bsPrefix}-backdrop`, backdropClassName)
  }), [backdropClassName, bsPrefix]);
  const renderDialog = (dialogProps) => (0, import_jsx_runtime8.jsx)("div", {
    ...dialogProps,
    ...props,
    className: (0, import_classnames7.default)(className, responsive ? `${bsPrefix}-${responsive}` : bsPrefix, `${bsPrefix}-${placement}`),
    "aria-labelledby": ariaLabelledby,
    children
  });
  return (0, import_jsx_runtime10.jsxs)(import_jsx_runtime9.Fragment, {
    children: [!showOffcanvas && (responsive || renderStaticNode) && renderDialog({}), (0, import_jsx_runtime8.jsx)(ModalContext_default.Provider, {
      value: modalContext,
      children: (0, import_jsx_runtime8.jsx)(Modal_default, {
        show: showOffcanvas,
        ref,
        backdrop,
        container,
        keyboard,
        autoFocus,
        enforceFocus: enforceFocus && !scroll,
        restoreFocus,
        restoreFocusOptions,
        onEscapeKeyDown,
        onShow,
        onHide: handleHide,
        onEnter: handleEnter,
        onEntering,
        onEntered,
        onExit,
        onExiting,
        onExited: handleExited,
        manager: getModalManager(),
        transition: DialogTransition,
        backdropTransition: BackdropTransition,
        renderBackdrop,
        renderDialog
      })
    })]
  });
});
Offcanvas.displayName = "Offcanvas";
var Offcanvas_default = Object.assign(Offcanvas, {
  Body: OffcanvasBody_default,
  Header: OffcanvasHeader_default,
  Title: OffcanvasTitle_default
});

// node_modules/react-bootstrap/esm/NavbarOffcanvas.js
var import_jsx_runtime11 = __toESM(require_jsx_runtime());
var NavbarOffcanvas = React9.forwardRef(({
  onHide,
  ...props
}, ref) => {
  const context = (0, import_react6.useContext)(NavbarContext_default);
  const handleHide = useEventCallback(() => {
    context == null || context.onToggle == null || context.onToggle();
    onHide == null || onHide();
  });
  return (0, import_jsx_runtime11.jsx)(Offcanvas_default, {
    ref,
    show: !!(context != null && context.expanded),
    ...props,
    renderStaticNode: true,
    onHide: handleHide
  });
});
NavbarOffcanvas.displayName = "NavbarOffcanvas";
var NavbarOffcanvas_default = NavbarOffcanvas;

// node_modules/react-bootstrap/esm/NavbarText.js
var React10 = __toESM(require_react());
var import_classnames8 = __toESM(require_classnames());
var import_jsx_runtime12 = __toESM(require_jsx_runtime());
var NavbarText = React10.forwardRef(({
  className,
  bsPrefix,
  as: Component = "span",
  ...props
}, ref) => {
  bsPrefix = useBootstrapPrefix(bsPrefix, "navbar-text");
  return (0, import_jsx_runtime12.jsx)(Component, {
    ref,
    className: (0, import_classnames8.default)(className, bsPrefix),
    ...props
  });
});
NavbarText.displayName = "NavbarText";
var NavbarText_default = NavbarText;

// node_modules/react-bootstrap/esm/Navbar.js
var import_jsx_runtime13 = __toESM(require_jsx_runtime());
var Navbar = React11.forwardRef((props, ref) => {
  const {
    bsPrefix: initialBsPrefix,
    expand = true,
    variant = "light",
    bg,
    fixed,
    sticky,
    className,
    // Need to define the default "as" during prop destructuring to be compatible with styled-components github.com/react-bootstrap/react-bootstrap/issues/3595
    as: Component = "nav",
    expanded,
    onToggle,
    onSelect,
    collapseOnSelect = false,
    ...controlledProps
  } = useUncontrolled(props, {
    expanded: "onToggle"
  });
  const bsPrefix = useBootstrapPrefix(initialBsPrefix, "navbar");
  const handleCollapse = (0, import_react7.useCallback)((...args) => {
    onSelect == null || onSelect(...args);
    if (collapseOnSelect && expanded) {
      onToggle == null || onToggle(false);
    }
  }, [onSelect, collapseOnSelect, expanded, onToggle]);
  if (controlledProps.role === void 0 && Component !== "nav") {
    controlledProps.role = "navigation";
  }
  let expandClass = `${bsPrefix}-expand`;
  if (typeof expand === "string") expandClass = `${expandClass}-${expand}`;
  const navbarContext = (0, import_react7.useMemo)(() => ({
    onToggle: () => onToggle == null ? void 0 : onToggle(!expanded),
    bsPrefix,
    expanded: !!expanded,
    expand
  }), [bsPrefix, expanded, expand, onToggle]);
  return (0, import_jsx_runtime13.jsx)(NavbarContext_default.Provider, {
    value: navbarContext,
    children: (0, import_jsx_runtime13.jsx)(SelectableContext_default.Provider, {
      value: handleCollapse,
      children: (0, import_jsx_runtime13.jsx)(Component, {
        ref,
        ...controlledProps,
        className: (0, import_classnames9.default)(className, bsPrefix, expand && expandClass, variant && `${bsPrefix}-${variant}`, bg && `bg-${bg}`, sticky && `sticky-${sticky}`, fixed && `fixed-${fixed}`)
      })
    })
  });
});
Navbar.displayName = "Navbar";
var Navbar_default = Object.assign(Navbar, {
  Brand: NavbarBrand_default,
  Collapse: NavbarCollapse_default,
  Offcanvas: NavbarOffcanvas_default,
  Text: NavbarText_default,
  Toggle: NavbarToggle_default
});
export {
  Navbar_default as default
};
//# sourceMappingURL=react-bootstrap_Navbar.js.map
