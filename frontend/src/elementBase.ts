import htmx from 'htmx.org';

export class FieldMapping {
    source: string
    targetSelector: string
    targetAttribute?: string
    isList?: boolean = false
}

// Base for elements with not shadow root
export class TemplElement extends HTMLElement {
    constructor(tplId: string) {
        super();
        const template = document.getElementById(tplId);
        if (!(template instanceof HTMLTemplateElement)) {
            return;
        }

        const templateContent = template.content;
        this.appendChild(document.importNode(templateContent, true));
    }
}

// Base for elements with shadow root
export class ShadowTemplElement extends HTMLElement {
    constructor(tplId: string, useGlobalStyles = false, fieldMappings: FieldMapping[] = []) {
        super();
        const template = document.getElementById(tplId);
        if (!(template instanceof HTMLTemplateElement)) {
            return;
        }

        const templateContent = template.content;
        const shadowRoot = this.attachShadow({ mode: "open" });
        shadowRoot.appendChild(document.importNode(templateContent, true));

        this.setupSlotLinks()
        this.copyDataFields(fieldMappings)

        if (useGlobalStyles) {
            const globalStyles = document.querySelectorAll('style'); // or any identifier
            globalStyles.forEach(style => {
                shadowRoot.appendChild(style.cloneNode(true));
            });
        }
    }

    connectedCallback() {
        // htmx does not auto-scan shadow roots; process this component root explicitly.
        if (htmx && typeof htmx.process === 'function') {
            htmx.process(this.shadowRoot);
        }
    }

    setupSlotLinks() {
        const slots = this.shadowRoot?.querySelectorAll<HTMLSlotElement>('slot');
        if (slots.length == 0) {
            return;
        }
        const slot = slots[0]
        slot.addEventListener('click', event => {
            const path = event.composedPath();
            const link = path.find(node => node instanceof HTMLAnchorElement);
            if (!(link instanceof HTMLAnchorElement)) {
                return;
            }

            const url = link.getAttribute('shadow-href-url') || link.getAttribute('href');
            if (!url || url === '#') {
                return;
            }

            const targetId = link.getAttribute('shadow-href-target')

            const target = this.shadowRoot?.querySelector<HTMLElement>(targetId);
            if (!target) {
                return;
            }

            event.preventDefault();
            event.stopPropagation();
            htmx.ajax('GET', url, target);
        });
    }

    copyDataFields(fieldMappings: FieldMapping[]) {
        const root = this.shadowRoot;
        const mappings = fieldMappings ?? [];

        if (!root || mappings.length === 0) {
            return;
        }

        for (const mapping of mappings) {
            if (!mapping?.source || !mapping?.targetSelector) {
                continue;
            }

            // Read source from host attributes, with a data-* fallback.
            let value = this.getAttribute(mapping.source);
            if (value === null && !mapping.source.startsWith('data-')) {
                value = this.getAttribute(`data-${mapping.source}`);
            }
            if (value === null) {
                continue;
            }

            const targets = root.querySelectorAll<HTMLElement>(mapping.targetSelector);
            targets.forEach(target => {
                if (mapping.targetAttribute) {
                    target.setAttribute(mapping.targetAttribute, value as string);
                } else {
                    target.textContent = value;
                }
            });
        }
    }
}

export function registerElementWithTemplate(
    elemId: string,
    elemCls: CustomElementConstructor,
    templateSource: string,
) {
    let template = document.getElementById(elemId);
    if (!(template instanceof HTMLTemplateElement)) {
        template = document.createElement('template');
        template.id = elemId;
        document.body.appendChild(template);
    }

    template.innerHTML = templateSource.trim();

    if (!customElements.get(elemId)) {
        customElements.define(elemId, elemCls);
    }
}

export function registerModal(elemId: string, templateSourece: string) { }