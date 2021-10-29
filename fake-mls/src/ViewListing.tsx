import React from 'react';
import { Card, Container, Grid } from '@material-ui/core';
import styled from '@emotion/styled'
import imgFakeMLS from './fake-mls.jpg';
import imgListing from './123mainst.jpg';

type InputEvent = React.ChangeEvent<HTMLInputElement>;

interface Data {
  address?: string,
  city?: string,
  state?: string,
  zip?: string,
};

interface Props {
  data: Data
};

function ViewListing({data}: Props) {
  return (
    <StyledContainer>
      <Grid container spacing={4}>
        <Grid item xs={5}>
          <Logo src={imgFakeMLS} alt="logo" />
        </Grid>
        <Grid item xs={6}>
          <h1>View Listing</h1>
          <StyledCard>
            <ListingPhoto>
              <img src={imgListing} alt="920 Fifth ave." />
            </ListingPhoto>
            <div>{data.address}</div>
            <div>
              {data.city}, {data.state} {data.zip}
            </div>
          </StyledCard>
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
const StyledCard = styled(Card)`
  font-size: 20px;
  padding: 30px;
  & > div {
    margin: 20px 0;
  }
`
const ListingPhoto = styled.div`
  display: flex;
  img {
    width: 100%;
    border: 1px solid #888;
  }
`

export default ViewListing;
