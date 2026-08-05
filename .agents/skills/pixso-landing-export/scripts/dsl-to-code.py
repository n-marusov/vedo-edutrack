#!/usr/bin/env python3
"""
Pixso DSL → React+Tailwind Code Generator

Takes the compact DSL JSON from get_node_dsl(simplify=True) and generates
React+Tailwind JSX code that faithfully reproduces the design.

Usage:
    python dsl-to-code.py <input.json> [--section-name <name>] [--output <file>]
    python dsl-to-code.py <input.json> --data-only   # Only extract data arrays (PLANS, FAQS, etc.)
    python dsl-to-code.py <input.json> --structure   # Only generate JSX structure (no data)

Input: compact DSL JSON file (from get_node_dsl with simplify=True)
Output: React+Tailwind JSX code to stdout or --output file
"""

import json
import os
import re
import sys
import argparse
from pathlib import Path

# ─── Variable resolution ───────────────────────────────────────────────

# Known CSS variable names mapped from Pixso variable names
# These are the variables defined in frontend/src/styles/pixso-variables.css
KNOWN_VARIABLES = {
    "background": "background",
    "foreground": "foreground",
    "muted": "muted",
    "muted-foreground": "muted-foreground",
    "card": "card",
    "border": "border",
    "primary-warm": "primary-warm",
    "secondary-calm": "secondary-calm",
    "accent-warm": "accent-warm",
    "accent-foreground": "accent-foreground",
    "primary-foreground": "primary-foreground",
    "primary": "primary",
    "primary-hover": "primary-hover",
    "destructive": "destructive",
    "ring": "ring",
    "warning": "warning",
    "success": "success",
    "card-hover": "card-hover",
}

# Lucide React icon name mapping based on Pixso vector names
ICON_MAP = {
    "hand": "Hand",
    "hand0": "Hand",
    "sparkles": "Sparkles",
    "arrowright": "ArrowRight",
    "arrowright0": "ArrowRight",
    "arrowright1": "ArrowRight",
    "arrowright2": "ArrowRight",
    "arrow-right": "ArrowRight",
    "rocket": "Rocket",
    "target": "Target",
    "target0": "Target",
    "clock": "Clock",
    "clock0": "Clock",
    "search": "Search",
    "unlink": "Unlink",
    "question": "HelpCircle",
    "alertcircle": "AlertCircle",
    "alertcircle0": "AlertCircle",
    "flaskconical": "FlaskConical",
    "folderkanban": "FolderKanban",
    "refreshcw": "RefreshCw",
    "check": "Check",
    "check0": "Check",
    "check1": "Check",
    "check2": "Check",
    "check3": "Check",
    "check4": "Check",
    "check5": "Check",
    "check6": "Check",
    "check7": "Check",
    "check8": "Check",
    "check9": "Check",
    "check10": "Check",
    "check11": "Check",
    "check12": "Check",
    "check13": "Check",
    "check14": "Check",
    "check15": "Check",
    "check16": "Check",
    "check17": "Check",
    "check18": "Check",
    "check19": "Check",
    "check20": "Check",
    "check21": "Check",
    "check22": "Check",
    "check23": "Check",
    "check24": "Check",
    "plus": "Plus",
    "plus0": "Plus",
    "plus1": "Plus",
    "plus2": "Plus",
    "plus3": "Plus",
    "minus": "Minus",
    "minus0": "Minus",
    "quote": "Quote",
    "quote0": "Quote",
    "quote1": "Quote",
    "compass": "Compass",
    "compass0": "Compass",
    "layers": "Layers",
    "link": "Link2",
    "link0": "Link2",
    "users": "Users",
    "users0": "Users",
    "users1": "Users",
    "eye": "Eye",
    "badgecheck": "BadgeCheck",
    "map": "Map",
    "map0": "Map",
    "navigation": "Navigation",
    "info": "Info",
    "briefcase": "Briefcase",
    "externallink": "ExternalLink",
    "send": "Send",
    "play": "Play",
    "send": "Send",
    "star": "Star",
    "search": "Search",
    "clock": "Clock",
    "eye": "Eye",
    "link": "Link2",
    "map": "Map",
    "badge-check": "BadgeCheck",
    "message-circle": "MessageCircle",
}

# ─── Pixel → Tailwind value conversion ─────────────────────────────────

def px_to_tailwind(value: float) -> str:
    """Convert a pixel value to a Tailwind arbitrary value string."""
    if value == 0:
        return "0"
    # Round to nearest integer
    rounded = round(value)
    if rounded == value:
        return f"[{int(value)}px]"
    return f"[{value}px]"

def px_to_tailwind_gap(value: float) -> str:
    """Convert pixel gap to Tailwind gap class."""
    common = {4: "1", 8: "2", 12: "3", 16: "4", 20: "5", 24: "6",
              32: "8", 40: "10", 48: "12", 64: "16", 72: "18",
              80: "20", 96: "24"}
    if value in common:
        return f"gap-{common[value]}"
    return f"gap-{px_to_tailwind(value)}"

def px_to_tailwind_padding(value, axis=None):
    """Convert pixel padding to Tailwind padding class."""
    common = {4: "1", 8: "2", 12: "3", 16: "4", 20: "5", 24: "6",
              32: "8", 40: "10", 48: "12", 64: "16", 80: "20",
              100: "25", 140: "36", 160: "40"}
    if value in common:
        suffix = common[value]
    else:
        suffix = px_to_tailwind(value)

    prefix = {"top": "pt", "bottom": "pb", "left": "pl", "right": "pr",
              "x": "px", "y": "py"}.get(axis, "p")
    return f"{prefix}-{suffix}"

# ─── Fill → CSS var conversion ────────────────────────────────────────

def resolve_fill(fill, variables):
    """Convert a fill to CSS var() or direct color."""
    if isinstance(fill, str):
        # variable:3:11 format
        m = re.match(r'variable:(\d+:\d+)', fill)
        if m:
            var_id = m.group(1)
            if var_id in variables:
                var_name = variables[var_id]
                if var_name in KNOWN_VARIABLES:
                    return f"var(--{var_name})"
                return f"var(--{var_name})"
        return fill

    if isinstance(fill, dict):
        if fill.get("type") == "solid":
            color = fill.get("value", "")
            return color
        # Handle gradient
        if "GRADIENT" in fill.get("type", ""):
            stops = fill.get("stops", [])
            if stops:
                # Generate CSS gradient
                colors = []
                for s in stops:
                    c = s.get("color", {})
                    pos = s.get("position", 0)
                    r, g, b, a = c.get("r", 0), c.get("g", 0), c.get("b", 0), c.get("a", 1)
                    colors.append(f"rgba({r},{g},{b},{a}) {pos * 100}%")
                direction = fill.get("type", "GRADIENT_LINEAR")
                if "LINEAR" in direction:
                    return f"linear-gradient(to right, {', '.join(colors)})"
                if "RADIAL" in direction:
                    return f"radial-gradient(circle, {', '.join(colors)})"
        return str(fill)

    return str(fill)

def extract_fills(fills, variables):
    """Extract fill color(s) from a node's fills array."""
    if not fills:
        return None
    for f in fills:
        if f.get("type") == "SOLID":
            color = f.get("color", {})
            r, g, b, a = color.get("r", 0), color.get("g", 0), color.get("b", 0), color.get("a", 1)
            if a == 0:
                continue  # Skip fully transparent fills
            return resolve_fill(f, variables)
        if f.get("type") == "GRADIENT_LINEAR":
            return resolve_fill(f, variables)
        if f.get("type") == "GRADIENT_RADIAL":
            return resolve_fill(f, variables)
    return None

# ─── AutoLayout → Tailwind class generation ───────────────────────────

def auto_layout_classes(al, node_id=None, node_name=None):
    """Convert autoLayout to Tailwind classes."""
    classes = []
    if not al:
        return classes

    direction = al.get("direction", "VERTICAL")
    if direction == "HORIZONTAL":
        classes.append("flex-row")
    else:
        classes.append("flex-col")

    # Gap
    gap = al.get("gap")
    if gap and gap > 0:
        classes.append(px_to_tailwind_gap(gap))

    # Alignment
    align = al.get("align", "center")
    primary_align = al.get("primaryAlign", "")
    counter_align = al.get("counterAlign", "")
    item_hori = al.get("itemHoriAlign", "")
    item_vert = al.get("itemVertAlign", "")

    if direction == "HORIZONTAL":
        # Primary axis = horizontal (justify)
        # Counter axis = vertical (items)
        if align == "space-between" or primary_align == "space-between" or item_hori == 3:
            classes.append("justify-between")
        elif align == "center" or primary_align == "center" or item_hori == 1:
            classes.append("justify-center")
        elif primary_align == "flex-start" or item_hori == 0:
            pass  # default justify-start
        elif item_hori == 2:
            classes.append("justify-end")

        if align == "center" or counter_align == "center" or item_vert == 1:
            classes.append("items-center")
        elif counter_align == "flex-start" or item_vert == 0:
            classes.append("items-start")
        elif item_vert == 2:
            classes.append("items-end")
        elif counter_align == "flex-end":
            classes.append("items-end")
    else:
        # Vertical: primary = vertical (justify), counter = horizontal (items)
        if align == "center" or primary_align == "center" or item_vert == 1:
            classes.append("justify-center")
        elif primary_align == "flex-start" or item_vert == 0:
            pass
        elif item_vert == 2:
            classes.append("justify-end")

        if align == "center" or counter_align == "center" or item_hori == 1:
            classes.append("items-center")
        elif counter_align == "flex-start" or item_hori == 0:
            classes.append("items-start")
        elif item_hori == 2:
            classes.append("items-end")

    # Width sizing
    width_resize = al.get("widthResize")
    if width_resize == 2:  # FILL
        classes.append("w-full")
    elif width_resize == 1:  # HUG
        classes.append("w-fit")

    height_resize = al.get("heightResize")
    if height_resize == 2:  # FILL
        classes.append("h-full")

    # Padding
    padding = al.get("padding")
    if padding:
        if isinstance(padding, list) and len(padding) == 4:
            t, r, b, l = padding
            if t == r == b == l and t > 0:
                classes.append(px_to_tailwind_padding(t))
            else:
                if t > 0: classes.append(px_to_tailwind_padding(t, "y") if t == b else px_to_tailwind_padding(t, "top"))
                if b > 0 and b != t: classes.append(px_to_tailwind_padding(b, "bottom"))
                if l > 0: classes.append(px_to_tailwind_padding(l, "x") if l == r else px_to_tailwind_padding(l, "left"))
                if r > 0 and r != l: classes.append(px_to_tailwind_padding(r, "right"))
                # Simplify if y and x are the same
                if t == b and r == l and t > 0 and r > 0:
                    classes = [c for c in classes if "top" not in c and "bottom" not in c and "left" not in c and "right" not in c]
                    if t == r:
                        classes.append(px_to_tailwind_padding(t))
                    else:
                        classes.append(px_to_tailwind_padding(t, "y"))
                        classes.append(px_to_tailwind_padding(r, "x"))
        elif isinstance(padding, (int, float)):
            if padding > 0:
                classes.append(px_to_tailwind_padding(padding))

    return classes

# ─── Box sizing ────────────────────────────────────────────────────────

def box_classes(node, al):
    """Generate width/height classes from box dimensions."""
    classes = []
    box = node.get("box", {})
    w = box.get("w", 0)
    h = box.get("h", 0)

    width_resize = al.get("widthResize") if al else None
    height_resize = al.get("heightResize") if al else None

    if width_resize is None or width_resize == 0:  # FIXED
        if w > 0 and w < 2000:
            classes.append(f"w-{px_to_tailwind(w)}")
    if height_resize is None or height_resize == 0:  # FIXED
        if h > 0 and h < 2000:
            classes.append(f"h-{px_to_tailwind(h)}")

    return classes

# ─── Radius conversion ─────────────────────────────────────────────────

def radius_class(radius):
    """Convert radius to Tailwind rounded class."""
    if not radius or radius == 0:
        return ""
    common = {4: "rounded-sm", 8: "rounded", 12: "rounded-xl",
              16: "rounded-2xl", 20: "rounded-[20px]", 24: "rounded-[24px]",
              28: "rounded-[28px]", 32: "rounded-[32px]", 40: "rounded-[40px]",
              999: "rounded-full"}
    if radius in common:
        return common[radius]
    return f"rounded-{px_to_tailwind(radius)}"

# ─── Shadow conversion ────────────────────────────────────────────────

def shadow_class(effects):
    """Convert Pixso effects to Tailwind shadow class."""
    if not effects:
        return ""
    for e in effects:
        if e.get("type") == "DROP_SHADOW" and e.get("visible", True):
            color = e.get("color", {})
            r, g, b, a = color.get("r", 0), color.get("g", 0), color.get("b", 0), color.get("a", 1)
            offset_x = e.get("offset", {}).get("x", 0)
            offset_y = e.get("offset", {}).get("y", 0)
            radius = e.get("radius", 0)
            spread = e.get("spread", 0)
            return f"shadow-[{offset_x}px_{offset_y}px_{radius}px_{spread}px_rgba({r},{g},{b},{a})]"
    return ""

# ─── Font mapping ─────────────────────────────────────────────────────

def font_classes(font_data):
    """Convert Pixso font properties to Tailwind classes."""
    classes = []
    if not font_data:
        return classes

    family = font_data.get("fontFamily", "")
    size = font_data.get("fontSize", 16)
    weight = font_data.get("fontWeight", 400)
    style = font_data.get("fontStyle", "")
    line_height = font_data.get("lineHeightNumber", None)
    line_height_unit = font_data.get("lineHeightUnit", None)
    letter_spacing = font_data.get("letterSpacingNumber", None)
    text_align = font_data.get("textAlignHorizontal", None)

    # Font family
    if family:
        classes.append(f"font-['{family}']")

    # Font size
    common_sizes = {11: "text-[11px]", 12: "text-xs", 13: "text-[13px]",
                    14: "text-sm", 15: "text-[15px]", 16: "text-base",
                    17: "text-[17px]", 18: "text-lg", 20: "text-xl",
                    22: "text-[22px]", 24: "text-2xl", 28: "text-3xl",
                    32: "text-4xl", 48: "text-5xl", 64: "text-7xl", 72: "text-8xl"}
    if size in common_sizes:
        classes.append(common_sizes[size])
    else:
        classes.append(f"text-{px_to_tailwind(size)}")

    # Font weight
    weight_map = {400: "font-normal", 500: "font-medium", 600: "font-semibold",
                  700: "font-bold", 800: "font-extrabold", 900: "font-black"}
    if weight in weight_map:
        classes.append(weight_map[weight])
    elif "Bold" in style:
        classes.append("font-bold")
    elif "Medium" in style:
        classes.append("font-medium")
    elif "Semi" in style:
        classes.append("font-semibold")
    elif "Extra" in style:
        classes.append("font-extrabold")
    else:
        classes.append("font-normal")

    # Line height
    if line_height and line_height_unit:
        if line_height_unit == "PERCENT":
            pct = int(line_height)
            if pct in [100, 110]:
                classes.append("leading-none")
            elif pct == 130:
                classes.append("leading-tight")
            elif pct == 140:
                classes.append("leading-snug")
            elif pct == 150:
                classes.append("leading-normal")
            elif pct == 160:
                classes.append("leading-relaxed")
            elif pct == 170:
                classes.append("leading-loose")
            else:
                classes.append(f"leading-[{pct}%]")
        elif line_height_unit == "PIXELS":
            classes.append(f"leading-{px_to_tailwind(line_height)}")

    # Letter spacing
    if letter_spacing and letter_spacing != 0:
        classes.append(f"tracking-{px_to_tailwind(letter_spacing)}")

    # Text align
    if text_align and text_align != "left":
        classes.append(f"text-{text_align}")

    return classes

# ─── Icon name resolution ─────────────────────────────────────────────

def resolve_icon_name(svg_sha: str) -> str:
    """Resolve SVG sha to lucide-react icon name."""
    name = svg_sha
    # Remove extension
    if '.' in name:
        name = name.rsplit('.', 1)[0]
    # Remove trailing numbers
    name = re.sub(r'\d+$', '', name)
    # Remove leading numbers
    name = re.sub(r'^\d+', '', name)
    return ICON_MAP.get(name, name.title())

# ─── Code generation ──────────────────────────────────────────────────

INDENT = "  "

def generate_jsx(node, variables, depth=0, parent_al=None):
    """Recursively generate JSX from a DSL node."""
    indent = INDENT * depth
    node_type = node.get("type", "")
    node_name = node.get("name", "")
    al = node.get("autoLayout", {})
    children = node.get("children", [])

    if node_type == "TEXT" or node_type == "PARAGRAPH":
        return _generate_text(node, variables, depth)

    if node_type in ("VECTOR", "ELLIPSE", "PATH"):
        return _generate_vector(node, variables, depth)

    if node_type == "INSTANCE":
        return _generate_instance(node, variables, depth)

    if node_type == "FRAME" or node_type == "GROUP":
        return _generate_frame(node, variables, depth)

    # Fallback
    return f"{indent}<div>{node_name}</div>\n"


def _generate_text(node, variables, depth):
    """Generate JSX for a text node."""
    indent = INDENT * depth
    text_content = node.get("nodeText", node.get("text", {}).get("content", ""))
    if not text_content:
        return ""

    fills = node.get("fills") or node.get("fillPaints", [])
    fill_color = extract_fills(fills, variables) if fills else None

    font_data = {
        "fontFamily": node.get("fontFamily"),
        "fontSize": node.get("fontSize"),
        "fontWeight": node.get("fontWeight"),
        "fontStyle": node.get("fontStyle"),
        "lineHeightNumber": node.get("lineHeightNumber"),
        "lineHeightUnit": node.get("lineHeightUnit"),
        "letterSpacingNumber": node.get("letterSpacingNumber"),
        "textAlignHorizontal": node.get("textAlignHorizontal"),
    }
    fclasses = font_classes(font_data)

    # Text auto resize
    auto_resize = node.get("textAutoResize", "")
    if auto_resize == "HEIGHT":
        pass  # No width constraint needed

    classes = []
    if fill_color:
        classes.append(f"text-[{fill_color}]")
    classes.extend(fclasses)

    # Check if text has child spans
    child_spans = []
    for child in node.get("childNode", []):
        if child.get("type") == "SPAN" or child.get("type") == "PARAGRAPH":
            span_text = child.get("nodeText", "")
            if span_text:
                child_spans.append(span_text)

    if child_spans:
        # Multi-paragraph text
        result = ""
        for span in child_spans:
            cls = " ".join(classes)
            result += f"{indent}<p className=\"{cls}\">{span}</p>\n"
        return result
    else:
        cls = " ".join(classes)
        if len(text_content) > 80:
            return f'{indent}<p className="{cls}">{text_content}</p>\n'
        return f'{indent}<p className="{cls}">{text_content}</p>\n'


def _generate_vector(node, variables, depth):
    """Generate JSX for a vector/icon node."""
    indent = INDENT * depth
    svg_sha = node.get("svgSha", "")
    icon_name = resolve_icon_name(svg_sha)

    fills = node.get("fills") or node.get("fillPaints", [])
    fill_color = extract_fills(fills, variables) if fills else None

    box = node.get("box", {})
    w = box.get("w", 24)
    h = box.get("h", 24)

    if icon_name:
        color_attr = f' strokeWidth={{1.5}} className="text-[{fill_color}]"' if fill_color else ' strokeWidth={{1.5}}'
        return f'{indent}<{icon_name} size={{{int(w)}}}{color_attr} />\n'

    # Fallback: inline SVG placeholder
    return f'{indent}<span className="inline-block w-{px_to_tailwind(w)} h-{px_to_tailwind(h)} bg-current" />\n'


def _generate_instance(node, variables, depth):
    """Generate JSX for a component instance."""
    indent = INDENT * depth
    comp_ref = node.get("componentRef", "")
    comp_name = node.get("name", "Component")
    overrides = node.get("override", {})
    children = node.get("children", [])

    # Known component patterns
    if "Section Header" in comp_name:
        badge = ""
        title = ""
        subtitle = ""
        for child in children:
            if child.get("type") == "INSTANCE":
                for sub in child.get("children", []):
                    if sub.get("type") == "TEXT":
                        badge = sub.get("text", {}).get("content", "") or sub.get("override", {}).get("text", {}).get("content", "")
            if child.get("type") == "TEXT":
                text_content = child.get("text", {}).get("content", "") or child.get("override", {}).get("text", {}).get("content", "")
                if "H2" in child.get("name", ""):
                    title = text_content
                elif "Subtitle" in child.get("name", ""):
                    subtitle = text_content

        if badge:
            return f'{indent}<SectionHeader badge="{badge}" title="{title}" subtitle="{subtitle}" />\n'
        return f'{indent}<SectionHeader badge="{badge or ""}" title="{title}" subtitle="{subtitle}" />\n'

    if "CTA Button" in comp_name or "CTA" in comp_name:
        cta_text = ""
        for child in children:
            if child.get("type") == "TEXT":
                cta_text = child.get("text", {}).get("content", "") or child.get("override", {}).get("text", {}).get("content", "")
        if cta_text:
            return f'{indent}<OrangeButton>{cta_text}</OrangeButton>\n'
        return f'{indent}<OrangeButton>{comp_name}</OrangeButton>\n'

    if "Eyebrow" in comp_name:
        text = ""
        for child in children:
            if child.get("type") == "TEXT":
                text = child.get("text", {}).get("content", "") or child.get("override", {}).get("text", {}).get("content", "")
        return f'{indent}<SectionBadge>{text}</SectionBadge>\n'

    # Generic instance - generate as a div with children
    result = f"{indent}<div>\n"
    for child in children:
        result += generate_jsx(child, variables, depth + 1)
    result += f"{indent}</div>\n"
    return result


def _generate_frame(node, variables, depth):
    """Generate JSX for a frame/group node."""
    indent = INDENT * depth
    node_name = node.get("name", "")
    al = node.get("autoLayout", {})
    box = node.get("box", {})
    children = node.get("children", [])

    is_section = depth == 0
    classes = []

    # AutoLayout
    al_classes = auto_layout_classes(al)
    classes.extend(al_classes)

    # Width/Height from box
    bx_classes = box_classes(node, al)
    classes.extend(bx_classes)

    # Radius
    radius = node.get("radius") or node.get("cornerRadius")
    if radius:
        rclass = radius_class(radius)
        if rclass:
            classes.append(rclass)

    # Fill
    fills = node.get("fills") or node.get("fillPaints", [])
    if fills:
        fill_color = extract_fills(fills, variables)
        if fill_color and not any("bg-" in c for c in classes):
            classes.append(f"bg-[{fill_color}]")

    # Stroke/border
    strokes = node.get("strokes") or node.get("strokePaints", [])
    stroke_weight = node.get("strokeWeight") or node.get("borderWeight")
    if strokes:
        stroke_color = extract_fills(strokes, variables)
        if stroke_color:
            classes.append(f"border border-[{stroke_color}]")
    elif stroke_weight:
        classes.append("border")

    # Effects/shadow
    effects = node.get("effects") or node.get("effectiveEffects", [])
    if effects:
        shadow = shadow_class(effects)
        if shadow:
            classes.append(shadow)

    # If no children, just render a div
    if not children:
        # Check if there are child nodes in childNode
        child_nodes = node.get("childNode", [])
        if child_nodes:
            cls = " ".join(classes)
            result = f"{indent}<div className=\"{cls}\">\n"
            for child in child_nodes:
                result += generate_jsx(child, variables, depth + 1)
            result += f"{indent}</div>\n"
            return result
        cls = " ".join(classes)
        return f"{indent}<div className=\"{cls}\" />\n"

    # Has children
    cls = " ".join(classes)
    result = f"{indent}<div className=\"{cls}\">\n"
    for child in children:
        result += generate_jsx(child, variables, depth + 1)
    result += f"{indent}</div>\n"
    return result


# ─── Data extraction (for data arrays) ─────────────────────────────────

def extract_data_arrays(node, variables, arrays=None, path=""):
    """Extract structured data arrays from a node tree.
    
    Returns dict with keys like 'plans', 'testimonials', 'faqs', 'principles', 'problems'
    """
    if arrays is None:
        arrays = {}
    
    node_type = node.get("type", "")
    node_name = node.get("name", "")
    children = node.get("children", [])
    child_nodes = node.get("childNode", children)
    
    # Detect section types by name patterns
    name_lower = node_name.lower()
    
    # Pricing plans
    if "plan" in name_lower and node_type == "FRAME":
        plan_data = extract_plan(node, variables)
        if plan_data:
            arrays.setdefault("plans", []).append(plan_data)
    
    # Testimonials
    if "testimonial" in name_lower and node_type == "FRAME":
        testimonial = extract_testimonial(node, variables)
        if testimonial:
            arrays.setdefault("testimonials", []).append(testimonial)
    
    # FAQ items
    if "faq" in name_lower and node_type == "FRAME":
        faq = extract_faq(node, variables)
        if faq:
            arrays.setdefault("faqs", []).append(faq)
    
    # Principles
    if "principle" in name_lower and node_type == "FRAME":
        principle = extract_principle(node, variables)
        if principle:
            arrays.setdefault("principles", []).append(principle)
    
    # Problem cards
    if "card" in name_lower and node_type == "FRAME" and "plan" not in name_lower:
        problem = extract_problem(node, variables)
        if problem:
            arrays.setdefault("problems", []).append(problem)
    
    # Benefits
    if "benefit" in name_lower and node_type == "FRAME":
        benefit = extract_benefit(node, variables)
        if benefit:
            arrays.setdefault("benefits", []).append(benefit)
    
    # Metrics
    if "metric" in name_lower and node_type == "FRAME":
        metric = extract_metric(node, variables)
        if metric:
            arrays.setdefault("metrics", []).append(metric)
    
    # Recurse
    for child in child_nodes:
        child_path = f"{path}/{child.get('name', child.get('id', ''))}"
        extract_data_arrays(child, variables, arrays, child_path)
    
    return arrays


def extract_plan(node, variables):
    """Extract pricing plan data."""
    data = {}
    children = node.get("children", [])
    child_nodes = node.get("childNode", children)
    
    for child in child_nodes:
        name = child.get("name", "").lower()
        ctype = child.get("type", "")
        
        if "title" in name or "name" in name:
            # Find plan name text
            for sub in child.get("children", child.get("childNode", [])):
                if sub.get("type") in ("TEXT", "PARAGRAPH"):
                    data["name"] = sub.get("nodeText", sub.get("text", {}).get("content", ""))
        
        if "badge" in name:
            for sub in child.get("children", child.get("childNode", [])):
                if sub.get("type") in ("TEXT", "PARAGRAPH"):
                    data["badge"] = sub.get("nodeText", sub.get("text", {}).get("content", ""))
        
        if "price" in name:
            texts = []
            for sub in child.get("children", child.get("childNode", [])):
                if sub.get("type") in ("TEXT", "PARAGRAPH"):
                    t = sub.get("nodeText", sub.get("text", {}).get("content", ""))
                    texts.append(t)
            if len(texts) >= 2:
                data["price"] = texts[0]
                data["period"] = texts[1]
        
        if "annual" in name:
            texts = []
            for sub in child.get("children", child.get("childNode", [])):
                if sub.get("type") in ("TEXT", "PARAGRAPH"):
                    texts.append(sub.get("nodeText", sub.get("text", {}).get("content", "")))
            if len(texts) >= 1:
                data["annual_price"] = texts[0]
            if len(texts) >= 2:
                data["annual_old"] = texts[1]
            if len(texts) >= 3:
                data["annual_savings"] = texts[2]
        
        if "description" in name or "desc" in name:
            data["description"] = child.get("nodeText", child.get("text", {}).get("content", ""))
        
        if "features" in name or "feature" in name:
            features = []
            for sub in child.get("children", child.get("childNode", [])):
                if sub.get("type") == "FRAME":
                    feature_text = ""
                    is_highlighted = False
                    for item in sub.get("children", sub.get("childNode", [])):
                        if item.get("type") in ("TEXT", "PARAGRAPH"):
                            feature_text = item.get("nodeText", item.get("text", {}).get("content", ""))
                    fills = sub.get("fills") or sub.get("fillPaints", [])
                    if fills:
                        fill_color = extract_fills(fills, variables)
                        if fill_color and "muted" in fill_color:
                            is_highlighted = True
                    if feature_text:
                        features.append({"text": feature_text, "highlight": is_highlighted})
            if features:
                data["features"] = features
    
    # Find CTA button
    for child in child_nodes:
        if child.get("type") == "INSTANCE" and ("CTA" in child.get("name", "") or "footer" in child.get("name", "").lower()):
            for sub in child.get("children", []):
                if sub.get("type") == "TEXT":
                    cta_text = sub.get("override", {}).get("text", {}).get("content", "") or sub.get("text", {}).get("content", "")
                    if cta_text:
                        data["cta"] = cta_text
    
    if data.get("name") and data.get("price"):
        return data
    return None


def extract_testimonial(node, variables):
    """Extract testimonial data."""
    data = {}
    children = node.get("children", [])
    child_nodes = node.get("childNode", children)
    
    for child in child_nodes:
        ctype = child.get("type", "")
        name = child.get("name", "").lower()
        
        if ctype in ("TEXT", "PARAGRAPH") and "quote" not in name:
            text = child.get("nodeText", child.get("text", {}).get("content", ""))
            if text and text.startswith("«"):
                data["text"] = text
            elif text and not data.get("text"):
                data["text"] = text
        
        if "author" in name:
            info = child.get("children", child.get("childNode", []))
            if len(info) >= 2:
                data["author"] = info[0].get("nodeText", info[0].get("text", {}).get("content", ""))
                data["role"] = info[1].get("nodeText", info[1].get("text", {}).get("content", ""))
            elif len(info) >= 1:
                data["author"] = info[0].get("nodeText", info[0].get("text", {}).get("content", ""))
        
        if "result" in name:
            for sub in child.get("children", child.get("childNode", [])):
                if sub.get("type") in ("TEXT", "PARAGRAPH"):
                    data["result"] = sub.get("nodeText", sub.get("text", {}).get("content", ""))
    
    return data if data.get("text") else None


def extract_faq(node, variables):
    """Extract FAQ item data."""
    data = {}
    children = node.get("children", [])
    child_nodes = node.get("childNode", children)
    
    for child in child_nodes:
        ctype = child.get("type", "")
        if ctype in ("TEXT", "PARAGRAPH"):
            text = child.get("nodeText", child.get("text", {}).get("content", ""))
            if child.get("visible") is False:
                data["answer"] = text
            elif not data.get("question"):
                data["question"] = text
    
    return data if data.get("question") and data.get("answer") else None


def extract_principle(node, variables):
    """Extract principle card data."""
    data = {}
    children = node.get("children", [])
    child_nodes = node.get("childNode", children)
    
    texts = []
    for child in child_nodes:
        ctype = child.get("type", "")
        name = child.get("name", "").lower()
        
        if ctype in ("TEXT", "PARAGRAPH"):
            text = child.get("nodeText", child.get("text", {}).get("content", ""))
            texts.append(text)
        elif "link" in name or child.get("link"):
            data["link"] = child.get("nodeText", child.get("text", {}).get("content", ""))
    
    # First text = title, rest = descriptions
    if len(texts) >= 1:
        data["title"] = texts[0]
    if len(texts) >= 2:
        data["subtitle"] = texts[1]
    if len(texts) >= 3:
        data["desc"] = texts[2]
    
    return data if data.get("title") else None


def extract_problem(node, variables):
    """Extract problem card data."""
    data = {}
    children = node.get("children", [])
    child_nodes = node.get("childNode", children)
    
    texts = []
    for child in child_nodes:
        ctype = child.get("type", "")
        if ctype in ("TEXT", "PARAGRAPH"):
            text = child.get("nodeText", child.get("text", {}).get("content", ""))
            texts.append(text)
    
    if len(texts) >= 1:
        data["title"] = texts[0]
    if len(texts) >= 2:
        data["desc"] = texts[1]
    
    return data if data.get("title") else None


def extract_benefit(node, variables):
    """Extract benefit card data."""
    data = {}
    children = node.get("children", [])
    child_nodes = node.get("childNode", children)
    
    texts = []
    for child in child_nodes:
        ctype = child.get("type", "")
        if ctype in ("TEXT", "PARAGRAPH"):
            text = child.get("nodeText", child.get("text", {}).get("content", ""))
            texts.append(text)
    
    # Detect "instead" text (starts with "Вместо:")
    instead = None
    desc_texts = []
    for t in texts:
        if t.startswith("Вместо:") or t.startswith("Вместо"):
            instead = t
        else:
            desc_texts.append(t)
    
    if len(desc_texts) >= 1:
        data["title"] = desc_texts[0]
    if len(desc_texts) >= 2:
        data["desc"] = desc_texts[1]
    if instead:
        data["instead"] = instead
    
    return data if data.get("title") else None


def extract_metric(node, variables):
    """Extract metric data."""
    data = {}
    children = node.get("children", [])
    child_nodes = node.get("childNode", children)
    
    for child in child_nodes:
        ctype = child.get("type", "")
        if ctype in ("TEXT", "PARAGRAPH"):
            text = child.get("nodeText", child.get("text", {}).get("content", ""))
            if any(c.isdigit() for c in text) and "+" in text:
                data["value"] = text
            else:
                data["label"] = text
    
    return data if data.get("value") else None


# ─── Section header detection ─────────────────────────────────────────

def detect_section_header(node):
    """Check if a node is a Section Header instance and extract props."""
    if node.get("type") != "INSTANCE":
        return None
    if "Section Header" not in node.get("name", ""):
        return None
    
    badge = ""
    title = ""
    subtitle = ""
    
    for child in node.get("children", []):
        if child.get("type") == "INSTANCE" and "Eyebrow" in child.get("name", ""):
            for sub in child.get("children", []):
                if sub.get("type") == "TEXT":
                    badge = sub.get("text", {}).get("content", "") or sub.get("override", {}).get("text", {}).get("content", "")
        if child.get("type") == "TEXT":
            name = child.get("name", "")
            text = child.get("text", {}).get("content", "") or child.get("override", {}).get("text", {}).get("content", "")
            if "H2" in name:
                title = text
            elif "Subtitle" in name:
                subtitle = text
    
    return {"badge": badge, "title": title, "subtitle": subtitle}


# ─── Main entry point ─────────────────────────────────────────────────

def main():
    parser = argparse.ArgumentParser(description="Pixso DSL → React+Tailwind code generator")
    parser.add_argument("input", help="DSL JSON file from get_node_dsl(simplify=True)")
    parser.add_argument("--section-name", "-n", help="Section name for the generated component", default="Section")
    parser.add_argument("--output", "-o", help="Output file for generated code")
    parser.add_argument("--data-only", action="store_true", help="Only extract data arrays")
    parser.add_argument("--structure", action="store_true", help="Only generate JSX structure")
    parser.add_argument("--list-icons", action="store_true", help="List all icon references found")
    args = parser.parse_args()
    
    with open(args.input, 'r', encoding='utf-8') as f:
        data = json.load(f)
    
    # Extract variables - can be either dict (root variableMap) or list (refsIndex)
    refs = data.get("refsIndex", {})
    variables = {}
    
    vars_data = refs.get("variables", {})
    if isinstance(vars_data, list):
        for v in vars_data:
            if isinstance(v, dict):
                vid = v.get("id", "")
                vname = v.get("name", vid)
                if vid:
                    variables[vid] = vname
    elif isinstance(vars_data, dict):
        for vid, vinfo in vars_data.items():
            if isinstance(vinfo, dict):
                variables[vid] = vinfo.get("name", vid)
            else:
                variables[vid] = str(vinfo)
    
    # Also check roots variableMap
    for root in data.get("roots", []):
        vmap = root.get("variableMap", {})
        for vid, vinfo in vmap.items():
            if isinstance(vinfo, dict):
                variables[vid] = vinfo.get("name", vid)
    
    roots = data.get("roots", [])
    if not roots:
        print("No roots found in DSL data", file=sys.stderr)
        sys.exit(1)
    
    root = roots[0]
    
    if args.data_only:
        # Extract data arrays
        arrays = extract_data_arrays(root, variables)
        # Output as JSON
        print(json.dumps(arrays, ensure_ascii=False, indent=2))
        return
    
    if args.list_icons:
        icons = set()
        def collect_icons(node):
            if node.get("type") == "VECTOR":
                sha = node.get("svgSha", "")
                if sha:
                    icons.add(sha)
            for child in node.get("children", node.get("childNode", [])):
                collect_icons(child)
        collect_icons(root)
        for icon in sorted(icons):
            mapped = resolve_icon_name(icon)
            print(f"{icon} → {mapped}")
        return
    
    # Generate full code
    output = ""
    
    # Check for section header
    section_header = None
    for child in root.get("children", root.get("childNode", [])):
        sh = detect_section_header(child)
        if sh:
            section_header = sh
            break
    
    if section_header:
        output += f"// Section: {root.get('name', args.section_name)}\n"
        output += f"// Header: {section_header['badge']} | {section_header['title']}\n"
        if section_header['subtitle']:
            output += f"// Subtitle: {section_header['subtitle']}\n"
        output += "\n"
    
    # Extract data arrays for reference
    arrays = extract_data_arrays(root, variables)
    for key, items in arrays.items():
        output += f"// Extracted data: {len(items)} {key}\n"
        for item in items:
            output += f"//   - {item.get('name') or item.get('title') or item.get('author') or item.get('question', '')[:60]}\n"
    output += "\n"
    
    # Generate JSX structure
    output += generate_jsx(root, variables)
    
    if args.output:
        Path(args.output).write_text(output, encoding="utf-8")
    else:
        print(output)


if __name__ == "__main__":
    main()
