var langTools = ace.require("ace/ext/language_tools");

var editor = ace.edit("editor");
editor.setTheme("ace/theme/solarized_dark");
editor.getSession().setMode("ace/mode/yaml");
editor.setShowPrintMargin(false);
editor.setOptions({
  enableBasicAutocompletion: true,
  enableSnippets: true,
  enableLiveAutocompletion: false,

  autoScrollEditorIntoView: true,
  fontSize: "16px",
  maxLines: 100,
  wrap: true,
});

function valid() {
  let getYamlCodeValidationErrors = (code) => {
    var error = "";
    try {
      jsyaml.safeLoad(code);
    } catch (e) {
      error = e;
    }
    return error;
  };

  let code = editor.getValue();
  let error = getYamlCodeValidationErrors(code);
  if (error) {
    editor.getSession().setAnnotations([
      {
        row: error.mark.line,
        column: error.mark.column,
        text: error.reason,
        type: "error",
      },
    ]);

    return false;
  } else {
    editor.getSession().setAnnotations([]);

    return true;
  }
}

function bin2string(array) {
  var result = "";
  for (var i = 0; i < array.length; ++i) {
    result += String.fromCharCode(array[i]);
  }
  return result;
}

function reload() {
  fetch("/api/reload", {
    method: "POST",
  }).then(function (response) {
    if (response.ok) {
      return response.text();
    } else {
      toastError(
        "Error saving configuration file. Response: " + response.status
      );
    }
  })
  .then(function (text) {
    toastInfo(text);
  });
}

function save() {
  fetch("/api/config", {
    method: "POST",
    body: editor.getValue(),
  })
    .then(function (response) {
      if (response.ok) {
        toastInfo("Configuration saved");
      } else {
        toastError(
          "Error saving configuration file. Response: " + response.status
        );
      }
    })
    .catch(function (error) {
      toastError("Error saving configuration file: " + error.message);
    });
}

editor.commands.addCommand({
  name: "save",
  bindKey: { win: "Ctrl-S", mac: "Command-S" },
  exec: function (editor) {
    if (valid()) {
      save();
    } else {
      toastError("Check file format errors before saving");
    }
  },
  readOnly: false,
});

editor.on("change", () => {
  valid();
});

fetch("/api/config")
  .then(function (response) {
    if (response.ok) {
      return response.text();
    } else {
      toastError(
        "Error getting data from server. Response: " + response.status
      );
    }
  })
  .then(function (yaml) {
    editor.setValue(yaml);
  })
  .catch(function (error) {
    toastError("Error getting yaml from server: " + error.message);
  });
