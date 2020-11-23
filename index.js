const express = require("express");
const multer  = require('multer');
const app = express();
const port = 3000;
//var upload = multer({ dest: 'uploads/' })
const { exec } = require("child_process");


const storage = multer.diskStorage({
    destination: "./uploads",
    filename: function(req, file, cb) {
      // null as first argument means no error
      //cb(null, Date.now() + "-" + file.originalname);
      cb(null, Date.now().toString());
      console.log("Is This where the name is created")
    },
  });

app.get("/", (req, res) => {
    //onrequest();
    res.sendFile('/Users/samuelvarghese/Desktop/Development/Web/Servero/index.html')
    //res.send("Hello")
});

app.listen(port, () => {
    console.log("Listening at http://localhost:", port);
});
const uploadStorage = multer({ storage: storage })

app.post("/upload", uploadStorage.single("file"), (req, res) => {
    onrequest(req.file.filename)
    console.log(req.file)
    return res.status(200)
  })  
//command for gpu temps and cpu temps
function onrequest(filename) {
    var directoryUploads = '/Users/samuelvarghese/Desktop/Development/Web/Servero/uploads'
    var fileDirectory = directoryUploads+'/'+filename
    var fullCommand = "lp "+ fileDirectory
    console.log("This is executed " + fullCommand)
    //console.log("WELL this is printing" + fullDirectory)
    exec(fullCommand, (error, stdout, stderr) => {
        if (error) {
            console.log("error printer");
            console.log(error);
            return;
        }
        if (stderr) {
            console.log("error ", stderr);
            return;
        }
        console.log("Output:");
        console.log(stdout);
    })
}