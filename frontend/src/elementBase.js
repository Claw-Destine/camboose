import htmx from 'htmx.org';

// Base for elements with not shadow root
export class TemplElement extends HTMLElement {
    constructor(tplId) {
        super();
        let template = document.getElementById(tplId);
        let templateContent = template.content;
        this.appendChild(document.importNode(templateContent, true));
    }
}

// Base for elements with shadow root
export class ShadowTemplElement extends HTMLElement {
    constructor(tplId, useGlobalStyles = false) {
        super();
        let template = document.getElementById(tplId);
        let templateContent = template.content;
        const shadowRoot = this.attachShadow({ mode: "open" });
        shadowRoot.appendChild(document.importNode(templateContent, true));

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
}

export function registerElementWithTemplate(elemId, elemCls, templateSource) {
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