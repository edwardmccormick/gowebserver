function ChatMessage({msg, idx}) {
    
  
      return (
              <div
                key={msg.id || idx}
                className={
                  msg.who === 'System'
                    ? 'text-center bg-warning text-white rounded m-2 p-2 flex-item'
                    : msg.who === 'Alert'
                    ? 'text-center bg-danger text-white rounded m-2 p-2 flex-item'
                    : msg.who === 'Connection'
                    ? 'text-center bg-success text-white rounded m-2 p-2 flex-item'
                    : msg.who === 'Admin' || msg.who === 'Moderator' || msg.who === "AI Host"
                    ? 'align-self-center bg-info text-white rounded m-2 p-2 flex-item'
                    : msg.who == `${User.id}` || msg.who === 'Me'
                    ? 'align-self-end bg-primary text-white rounded m-2 p-2 flex-item'
                    : 'align-self-start bg-light text-black rounded m-2 p-2 flex-item'
                }
              >
                {msg.time} - {msg.who} - {msg.message}
              </div>
    )
}
export default ChatMessage