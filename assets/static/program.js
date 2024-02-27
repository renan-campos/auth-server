import {BackendClient} from './javascript/client.js'

// Define a UI structure (could be an object or a class)
class UI {
    constructor(backendClient) {
        this.backendClient = backendClient;

        this.KeyInput = createTextBox("KeyInput", "Enter key");
        var keyInput = this.KeyInput
        var listener = function() {
            keyInput.value = "";
            keyInput.removeEventListener("focus", listener)
        };
        this.KeyInput.addEventListener("focus", listener);

        this.SendButton = createButton(
            'SendButton',
            'Send',
            () => {
               this.backendClient.SendOtp(this.KeyInput.value, (resp) => {
                   if (!resp.ok) {
                       this.KeyInput.value = "Invalid OTP, please try again";
                   } else {
                       this.KeyInput.value = "OTP Successful!";
                   }
               });
               this.KeyInput.addEventListener("focus", listener);
            });
    }

    Render() {
        for (const element of [
            this.KeyInput,
            this.SendButton
        ]) {
            document.body.appendChild(element)
        }
    }
}

// Create an instance of the UI structure
const ui = new UI(new BackendClient());
ui.Render();


function createButton(id, textContent, clickHandler) {
    var button = document.createElement('button');
    button.id = id
    button.textContent = textContent;
    button.addEventListener('click', clickHandler)
    return button
}

function createTextBox(textBoxId, defaultText) {
    var textBox = document.createElement("input");
    textBox.setAttribute("type", "text");
    textBox.setAttribute("id", textBoxId);
    textBox.setAttribute("value", defaultText);

    return textBox
}
