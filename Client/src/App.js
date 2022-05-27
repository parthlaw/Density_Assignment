import logo from "./logo.svg";
import "./App.css";
import { useState } from "react";
import Game from "./Game";
function App() {
  const [start, setStart] = useState(false);
  const [name, setName] = useState("");
  const handleEnterClick = () => {
    setStart(true);
  };
  return (
    <>
      {start === false && (
        <div>
          <label>
            Name:
            <input
              type="text"
              name="name"
              onChange={(e) => {
                setName(e.target.value);
              }}
            />
          </label>
          <button onClick={handleEnterClick}>Enter</button>
        </div>
      )}
      {start === true && <Game name={name} />}
    </>
  );
}

export default App;
