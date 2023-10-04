
declare module 'turndown' {
    const content: any = turndown
    export = content
}


declare module '@toast-ui/editor' {
    class Editor {
        static factory(arg0: { el: Element | null; viewer: boolean; height: string; initialValue: any; theme: string }) {
            throw new Error("Method not implemented.")
        }
        constructor(option: any)
        getHTML: () => any
        getMarkdown: () => any
        removeHook: (hook: string) => any
        addHook: (hook: string, callback: Function) => any
        setHeight: (height: string) => any
        getHeight: () => any
        on: (event: string, callback: Function) => any
        focus: () => any
        blur: () => any
        setMarkdown: (markdown: string, cursorToEnd: boolean) => any
        insertToolbarItem: (indexInfo: Object, item: string | Object) => any
    }
    export = Editor
}
