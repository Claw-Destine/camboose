import htmx from 'htmx.org';

function cloneTemplateContent(templateSource) {
    const template = document.createElement('template');
    template.innerHTML = templateSource.trim();
    return template.content.cloneNode(true);
}

// Base for elements with not shadow root
export class TemplElement extends HTMLElement {
    constructor(templateSource) {
        super();
        this.appendChild(cloneTemplateContent(templateSource));
    }
}

// Base for elements with shadow root
export class ShadowTemplElement extends HTMLElement {
    constructor(templateSource, useGlobalStyles = false) {
        super();
        const shadowRoot = this.attachShadow({ mode: "open" });
        shadowRoot.appendChild(cloneTemplateContent(templateSource));
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