function convertISODateToLocal(dateString) {
  if (!dateString) {
    return '';
  }

  const date = new Date(dateString);
  if (Number.isNaN(date.getTime())) {
    return String(dateString);
  }

  return date.toLocaleTimeString();
}

function ChatMessage({ msg, currentUserId, otherUserId, otherUserName }) {
  const isCurrentUser = msg.who === currentUserId || msg.who === 'Me';
  const isSystem = msg.who === 'System';
  const isAlert = msg.who === 'Alert';
  const isConnection = msg.who === 'Connection';
  const isAI = msg.who === 'AI Introduction' || msg.who === 'AI Host';
  const showTrailingTime = msg.who === otherUserId || msg.who === otherUserName || (!isCurrentUser && !isSystem && !isAlert && !isConnection);

  const containerClassName =
    isSystem
      ? 'text-center bg-warning text-white rounded m-2 p-2 flex-item'
      : isAlert
      ? 'text-center bg-danger text-white rounded m-2 p-2 flex-item'
      : isConnection
      ? 'text-center bg-success text-white rounded m-2 p-2 flex-item'
      : isAI
      ? 'align-self-center bg-purple gradient-bg text-white rounded m-2 p-2 flex-item'
      : isCurrentUser
      ? 'd-flex flex-row justify-content-end align-items-center'
      : 'd-flex flex-row justify-content-start align-items-center';

  const bubbleClassName =
    isSystem
      ? 'text-center bg-warning text-white rounded m-2 p-2 flex-item'
      : isAlert
      ? 'text-center bg-danger text-white rounded m-2 p-2 flex-item'
      : isConnection
      ? 'text-center bg-success text-white rounded m-2 p-2 flex-item'
      : isAI
      ? 'align-self-center bg-purple gradient-bg text-white rounded m-2 p-1 flex-item'
      : isCurrentUser
      ? 'bg-primary text-white rounded m-2 p-2 flex-item'
      : 'bg-light text-black rounded m-2 p-2 flex-item';

  return (
    <div className={containerClassName}>
      {isCurrentUser ? <span className="text-body-tertiary fs-6">{convertISODateToLocal(msg.time)}</span> : null}
      <span className={bubbleClassName}>
        {msg.who} - {msg.message}
      </span>
      {showTrailingTime ? <span className="text-body-tertiary fs-6">{convertISODateToLocal(msg.time)}</span> : null}
    </div>
  );
}

export default ChatMessage;
