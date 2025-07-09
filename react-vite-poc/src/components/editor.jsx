import React, { useState, useRef, useEffect } from 'react';

// --- Script and Stylesheet Loader ---
// This effect will load the Quill library and its CSS from a CDN.
// This is necessary because the environment doesn't have direct package imports.
export const useQuillLoader = () => {
    const [isLoaded, setIsLoaded] = useState(false);

    useEffect(() => {
        const script = document.createElement('script');
        script.src = 'https://cdn.quilljs.com/1.3.6/quill.js';
        script.async = true;
        script.onload = () => {
            setIsLoaded(true);
        };

        const link = document.createElement('link');
        link.rel = 'stylesheet';
        link.href = 'https://cdn.quilljs.com/1.3.6/quill.snow.css';

        document.head.appendChild(link);
        document.body.appendChild(script);

        // Cleanup function to remove the script and stylesheet when the component unmounts
        return () => {
            document.head.removeChild(link);
            document.body.removeChild(script);
        };
    }, []); // Empty dependency array means this effect runs only once

    return isLoaded;
};


// --- QuillEditor Component ---
// This component wraps the Quill editor instance.
export const QuillEditor = ({ onContentChange, initialContent = null }) => {
    const editorRef = useRef(null); // Ref to the div where Quill will be mounted
    const quillInstanceRef = useRef(null); // Ref to the Quill instance itself

    useEffect(() => {
        // We check if the editor div is rendered and if Quill is available on the window object
        if (editorRef.current && window.Quill && !quillInstanceRef.current) {
            // Initialize Quill
            const quill = new window.Quill(editorRef.current, {
                theme: 'snow', // Use the 'snow' theme
                modules: {
                    toolbar: [
                        [{ 'header': [1, 2, 3, false] }],
                        ['bold', 'italic', 'underline', 'strike'],
                        [{ 'list': 'ordered'}, { 'list': 'bullet' }],
                        ['link', 'image'],
                        ['clean']
                    ],
                },
                placeholder: 'This is the part where you pretend to have something interesting to say...',
            });

            // Store the Quill instance in the ref
            quillInstanceRef.current = quill;
            
            // Set initial content if provided
            if (initialContent) {
                quill.setContents(initialContent);
            }
            
            // --- Event Listener for Text Change ---
            quill.on('text-change', (delta, oldDelta, source) => {
                if (source === 'user') {
                    const content = quill.getContents();
                    if (onContentChange) {
                        onContentChange(content);
                    }
                }
            });
        }
    }, [onContentChange, initialContent]); // Rerun effect if onContentChange or initialContent changes

    return (
        // The div element that will be replaced by the Quill editor
        <div ref={editorRef} style={{ height: '250px', backgroundColor: '#fff' }}></div>
    );
};

export default QuillEditor;