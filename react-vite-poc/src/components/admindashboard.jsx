import { useState, useEffect } from 'react';
import { Tabs, Tab, Table, Container, Row, Col, Card, Badge } from 'react-bootstrap';
import { formatDistance } from 'date-fns';

function AdminDashboard({ jwt }) {
  const [users, setUsers] = useState([]);
  const [chats, setChats] = useState([]);
  const [recentUpdates, setRecentUpdates] = useState({
    recent_profiles: [],
    recent_photos: []
  });
  const [loading, setLoading] = useState({
    users: true,
    chats: true,
    recent: true
  });
  const [error, setError] = useState({
    users: null,
    chats: null,
    recent: null
  });

  // Fetch users data
  useEffect(() => {
    if (!jwt) return;

    setLoading(prev => ({ ...prev, users: true }));
    fetch('http://localhost:8080/admin/users', {
      headers: {
        'Authorization': jwt
      }
    })
      .then(res => {
        if (!res.ok) {
          throw new Error(`HTTP error! Status: ${res.status}`);
        }
        return res.json();
      })
      .then(data => {
        setUsers(data);
        setLoading(prev => ({ ...prev, users: false }));
      })
      .catch(err => {
        console.error('Error fetching users:', err);
        setError(prev => ({ ...prev, users: err.message }));
        setLoading(prev => ({ ...prev, users: false }));
      });
  }, [jwt]);

  // Fetch chats data
  useEffect(() => {
    if (!jwt) return;

    setLoading(prev => ({ ...prev, chats: true }));
    fetch('http://localhost:8080/admin/chats', {
      headers: {
        'Authorization': jwt
      }
    })
      .then(res => {
        if (!res.ok) {
          throw new Error(`HTTP error! Status: ${res.status}`);
        }
        return res.json();
      })
      .then(data => {
        setChats(data);
        setLoading(prev => ({ ...prev, chats: false }));
      })
      .catch(err => {
        console.error('Error fetching chats:', err);
        setError(prev => ({ ...prev, chats: err.message }));
        setLoading(prev => ({ ...prev, chats: false }));
      });
  }, [jwt]);

  // Fetch recent updates
  useEffect(() => {
    if (!jwt) return;

    setLoading(prev => ({ ...prev, recent: true }));
    fetch('http://localhost:8080/admin/recent', {
      headers: {
        'Authorization': jwt
      }
    })
      .then(res => {
        if (!res.ok) {
          throw new Error(`HTTP error! Status: ${res.status}`);
        }
        return res.json();
      })
      .then(data => {
        setRecentUpdates(data);
        setLoading(prev => ({ ...prev, recent: false }));
      })
      .catch(err => {
        console.error('Error fetching recent updates:', err);
        setError(prev => ({ ...prev, recent: err.message }));
        setLoading(prev => ({ ...prev, recent: false }));
      });
  }, [jwt]);

  // Function to make a user admin
  const makeUserAdmin = async (userId, isAdmin) => {
    try {
      const response = await fetch('http://localhost:8080/admin/set-admin', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': jwt
        },
        body: JSON.stringify({ user_id: userId, is_admin: isAdmin })
      });

      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }

      // Update the user in the list
      const data = await response.json();
      setUsers(users.map(user => 
        user.id === userId ? { ...user, is_admin: isAdmin } : user
      ));
      
      alert(`User ${isAdmin ? 'promoted to' : 'demoted from'} admin successfully`);
    } catch (err) {
      console.error('Error updating admin status:', err);
      alert(`Failed to update admin status: ${err.message}`);
    }
  };

  const formatDate = (dateString) => {
    if (!dateString) return 'N/A';
    try {
      const date = new Date(dateString);
      return formatDistance(date, new Date(), { addSuffix: true });
    } catch (e) {
      return dateString;
    }
  };

  return (
    <Container className="mt-4 mb-4">
      <h1 className="mb-4">Admin Dashboard</h1>
      
      <Tabs defaultActiveKey="users" className="mb-4">
        <Tab eventKey="users" title="Users">
          {loading.users ? (
            <p>Loading users data...</p>
          ) : error.users ? (
            <p className="text-danger">Error: {error.users}</p>
          ) : (
            <Table striped bordered hover responsive>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Email</th>
                  <th>Name</th>
                  <th>Age</th>
                  <th>Admin</th>
                  <th>Created</th>
                  <th>Last Login</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {users.map(user => (
                  <tr key={user.id}>
                    <td>{user.id}</td>
                    <td>{user.email}</td>
                    <td>{user.user?.name || 'No profile'}</td>
                    <td>{user.user?.age || 'N/A'}</td>
                    <td>{user.is_admin ? 'Yes' : 'No'}</td>
                    <td>{formatDate(user.create_time)}</td>
                    <td>{formatDate(user.last_login)}</td>
                    <td>
                      <button 
                        className={`btn btn-sm ${user.is_admin ? 'btn-danger' : 'btn-primary'}`}
                        onClick={() => makeUserAdmin(user.id, !user.is_admin)}
                      >
                        {user.is_admin ? 'Remove Admin' : 'Make Admin'}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </Table>
          )}
        </Tab>
        
        <Tab eventKey="chats" title="Recent Chats">
          {loading.chats ? (
            <p>Loading chat data...</p>
          ) : error.chats ? (
            <p className="text-danger">Error: {error.chats}</p>
          ) : (
            <Table striped bordered hover responsive>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Match ID</th>
                  <th>User ID</th>
                  <th>Time</th>
                  <th>Message</th>
                </tr>
              </thead>
              <tbody>
                {chats.slice(0, 50).map(message => (
                  <tr key={message.id}>
                    <td>{message.id}</td>
                    <td>{message.match_id}</td>
                    <td>{message.who}</td>
                    <td>{formatDate(message.time)}</td>
                    <td>{message.message}</td>
                  </tr>
                ))}
              </tbody>
            </Table>
          )}
        </Tab>
        
        <Tab eventKey="recent" title="Recent Updates">
          {loading.recent ? (
            <p>Loading recent updates...</p>
          ) : error.recent ? (
            <p className="text-danger">Error: {error.recent}</p>
          ) : (
            <Row>
              <Col md={6}>
                <Card>
                  <Card.Header>
                    <h5>Recently Updated Profiles</h5>
                  </Card.Header>
                  <Card.Body>
                    <Table striped bordered hover responsive>
                      <thead>
                        <tr>
                          <th>ID</th>
                          <th>Name</th>
                          <th>Updated</th>
                        </tr>
                      </thead>
                      <tbody>
                        {recentUpdates.recent_profiles.map(profile => (
                          <tr key={profile.id}>
                            <td>{profile.id}</td>
                            <td>{profile.name}</td>
                            <td>{formatDate(profile.update_time)}</td>
                          </tr>
                        ))}
                      </tbody>
                    </Table>
                  </Card.Body>
                </Card>
              </Col>
              
              <Col md={6}>
                <Card>
                  <Card.Header>
                    <h5>Recently Uploaded Photos</h5>
                  </Card.Header>
                  <Card.Body>
                    <Table striped bordered hover responsive>
                      <thead>
                        <tr>
                          <th>ID</th>
                          <th>User ID</th>
                          <th>Caption</th>
                          <th>Created</th>
                        </tr>
                      </thead>
                      <tbody>
                        {recentUpdates.recent_photos.map(photo => (
                          <tr key={photo.id}>
                            <td>{photo.id}</td>
                            <td>{photo.person_id}</td>
                            <td>{photo.caption}</td>
                            <td>{formatDate(photo.create_time)}</td>
                          </tr>
                        ))}
                      </tbody>
                    </Table>
                  </Card.Body>
                </Card>
              </Col>
            </Row>
          )}
        </Tab>
      </Tabs>
    </Container>
  );
}

export default AdminDashboard;