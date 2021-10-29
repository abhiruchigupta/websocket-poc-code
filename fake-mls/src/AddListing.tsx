import React, { SyntheticEvent, useState } from 'react';
import { Button, Container, Grid, TextField } from '@material-ui/core';
import styled from '@emotion/styled'
import imgFakeMLS from './fake-mls.jpg';
import { useHistory } from 'react-router-dom';

type InputEvent = React.ChangeEvent<HTMLInputElement>;

interface Props {
  submit: (data: any) => void
}

function AddListing(props: Props) {
  const history = useHistory()
  const [address, setAddress] = useState('')
  const [city, setCity] = useState('')
  const [state, setState] = useState('')
  const [zip, setZip] = useState('')
  const handleSubmit = async (e: SyntheticEvent) => {
    e.preventDefault()
    const data = {
      address,
      city,
      state,
      zip
    }
    console.log("DATA1:", data)
    await props.submit(data)
    history.push('/view')
  }
  return (
    <StyledContainer>
      <Grid container spacing={4}>
        <Grid item xs={5}>
          <Logo src={imgFakeMLS} alt="logo" />
        </Grid>
        <Grid item xs={6}>
          <h1>Add a Listing</h1>
          <div>
            <form noValidate onSubmit={handleSubmit}>
              <TextField
                value={address}
                onChange={(e: InputEvent) => setAddress(e.target.value)}
                variant="standard"
                margin="normal"
                required
                fullWidth
                label="Address"
                id="address"
                name="address"
                autoComplete="address"
                autoFocus
              />
              <div>
                <CityTextField
                  value={city}
                  onChange={(e: InputEvent) => setCity(e.target.value)}                
                  variant="standard"
                  margin="normal"
                  required
                  label="City"
                  id="city"
                  name="city"
                  autoComplete="city"
                />
                <StateTextField                
                  value={state}
                  onChange={(e: InputEvent) => setState(e.target.value)}
                  variant="standard"
                  margin="normal"
                  required
                  label="State"
                  id="state"
                  name="state"
                  autoComplete="state"
                />
                <ZipTextField
                  value={zip}
                  onChange={(e: InputEvent) => setZip(e.target.value)}                
                  type="zipcode"
                  variant="standard"
                  margin="normal"
                  required
                  label="Zipcode"
                  id="zipcode"
                  name="zipcode"
                  autoComplete="zipcode"
                />
              </div>
              <SubmitButton
                type="submit"
                fullWidth
                variant="contained"
                color="primary"
              >Create</SubmitButton>
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
const CityTextField = styled(TextField)`
  margin-right: 20px;
  width: 300px;
`
const StateTextField = styled(TextField)`
  margin-right: 20px;
  width: 80px;
`
const ZipTextField = styled(TextField)`
  width: 120px;
`
const SubmitButton = styled(Button)`
  margin: 20px 0;
`;

export default AddListing;
