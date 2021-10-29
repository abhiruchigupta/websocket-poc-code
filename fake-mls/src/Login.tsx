import React, { SyntheticEvent, useState } from 'react';
import { Button, Container, Grid, TextField } from '@material-ui/core';
import styled from '@emotion/styled'
import imgFakeMLS from './fake-mls.jpg';
import { useHistory } from 'react-router-dom';

type InputEvent = React.ChangeEvent<HTMLInputElement>;

interface Props {
  submit: (data: any) => void
}

function Login(props: Props) {
  const history = useHistory()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')

  const handleSubmit = (e: SyntheticEvent) => {
    e.preventDefault()
    const data = {
      email
    }
    console.log('LOGIN DATA:', data)
    props.submit(data)
    history.push('/add')
  }

  return (
    <StyledContainer>
      <Grid container spacing={4}>
        <Grid item xs={5}>
          <Logo src={imgFakeMLS} alt="logo" />
        </Grid>
        <Grid item xs={6}>
          <h1>Welcome to Fake MLS</h1>
          <div>
            <form noValidate onSubmit={handleSubmit}>
              <TextField
                value={email}
                onChange={(e:InputEvent) => setEmail(e.target.value)}
                variant="outlined"
                margin="normal"
                required
                fullWidth
                label="Email"
                id="email"
                name="email"
                autoComplete="email"
                autoFocus
              />
              <TextField
                value={password}
                onChange={(e:InputEvent) => setPassword(e.target.value)}
                type="password"
                variant="outlined"
                margin="normal"
                required
                fullWidth
                label="Password"
                id="password"
                name="password"
                autoComplete="password"
              />
              <SubmitButton
                type="submit"
                fullWidth
                variant="contained"
                color="primary"
              >Login</SubmitButton>
            </form>
          </div>
        </Grid>
      </Grid>
    </StyledContainer>
  );
}

const StyledContainer = styled(Container)`
  background-color: white;
  margin-top: 50px;
`
const Logo = styled.img`
  width: 100%;
`
const SubmitButton = styled(Button)`
  margin: 20px 0;
`;

export default Login;
