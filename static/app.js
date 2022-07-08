function $(query) {
    return document.querySelector(query)
}

const root = ReactDOM.createRoot($("#app"))


let files;

const Link = ({text, href}) => {
    return <a href={href}>Link: <span className={"underline"}>{text}</span></a>
}

const Form = () => {

    const [dragActive, setDragActive] = React.useState(false)
    const [uploadLink, setUploadLink ] = React.useState("")
    const [uploadText, setUploadText ] = React.useState("")

    const handleDrag = (e) => {
        e.preventDefault()
        e.stopPropagation()

        if (e.type === "dragenter" || e.type === "dragover") {
            setDragActive(true)
            $("#form-header").innerText = "Drop file!"
        } else if (e.type === "dragleave") {
            setDragActive(false)
            $("#form-header").innerText = "Drag file here to upload"
        }
    }

    const handleDrop = async (e) => {
        e.preventDefault()
        e.stopPropagation()
        setDragActive(false)

        if (e.dataTransfer.files && e.dataTransfer.files[0]) {
            $("#form-header").innerText = "Uploading..."

            const formData = new FormData()
            formData.append("fileUpload", e.dataTransfer.files[0]);

            const res = await fetch("/api/upload", {
                method: "POST",
                mode: "no-cors",
                accept: "*",
                AccessControlAllowOrigin: "*",
                body: formData
            })

            const body = await res.json()

                setUploadLink(body.url)
                setUploadText(body.filename)
                $("#form-header").innerText = "Uploaded!"
        }
    }

    const styles = {
        form: "font-bold text-pink-200",
        dragArea: "min-h-screen min-w-screen flex flex-col items-center justify-center gap-2",
        headerText: "select-none",
        dragFileElement: "w-full h-full absolute rounded-md top-0 right-0 left-0 bottom-0 bg-pink-200/20",
    }

    return <form
        className={styles.form}
        encType="multipart/form-data"
        id="upload-form"
        onDragEnter={handleDrag}
        onSubmit={(e) => e.preventDefault()}>
        <div
            className={styles.dragArea}
        >
            <header id="form-header" className={styles.headerText}>Drag file here to upload</header>
            <input multiple={true} id="fileUpload" name="fileUpload" type="file" hidden />
            { dragActive && <div className={styles.dragFileElement} onDragEnter={handleDrag} onDragLeave={handleDrag} onDragOver={handleDrag} onDrop={handleDrop} /> }
            { uploadLink && <Link text={uploadText} href={uploadLink}/> }
        </div>
    </form>
}

root.render(
    <Form />
)