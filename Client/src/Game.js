import React, { useEffect, useRef, useState } from "react";

const Game = ({ name }) => {
  const sock = useRef();
  const [time, setTime] = useState(0);
  const [message, setMessage] = useState("");
  useEffect(() => {
    const ws = new WebSocket(`ws://localhost:8080/ws?name=${name}`);
    sock.current = ws;
    ws.onmessage = function (event) {
      const data = JSON.parse(event.data);
      console.log(data.action);
      if (data.action === "sync") {
        console.log(parseInt(data.message));
        setTime(parseInt(data.message));
      }
      if (data.action === "update_counter") {
        setMessage(data.message);
      }
    };
  }, [name]);
  React.useEffect(() => {
    console.log(time);
    setTimeout(() => setTime((time + 1) % 60), 1000);
  }, [time]);
  const handleClick = (mesg) => {
    const message = {
      action: "place_bet",
      message: mesg,
    };
    sock.current.send(JSON.stringify(message));
  };
  return (
    <>
      <div>Game</div>
      <span>{time}</span>
      <span>{message}</span>
      <div>
        <button
          onClick={() => {
            handleClick("Up");
          }}
        >
          UP
        </button>
        <button
          onClick={() => {
            handleClick("Down");
          }}
        >
          DOWN
        </button>
      </div>
    </>
  );
};

export default Game;
