Feature: List payments
    Clients should be able to retrieve details about a collection of payments

    Scenario: List payments without pagination (API defaults to 10 payments/page)
        Given there are payments with IDs:
            | 29bc7b20-f134-4aee-954d-36fb1daa8c13 |
            | ac6bd520-c5b9-4573-969a-99d46690aa26 |
            | 233f2a4b-631e-4f0c-a1dc-c33a3792e9b1 |
            | f2fdea8a-2358-4459-b126-eda71065d5c5 |
            | d73e8403-c6a3-49eb-b24b-8c441496a5fd |
            | 9c27337b-2c8a-40b1-abab-a51ec8d6561c |
            | 34d746ce-1f7d-4ea6-958d-28dda4b2c3cb |
            | 199aa4d7-1af9-4933-9948-44b98d8ba03e |
            | eee4cb9a-f9c1-4d51-95f3-0647fb174124 |
            | 80f43626-e895-447e-99b6-6b260627f3ca |
            | 23ffa122-44aa-4106-b036-877e61848d29 |
            | 8d05996d-aad6-4ff2-813a-91130a4188f7 |
            | 5fd8ed2b-407a-4089-ae25-b2fd70a83250 |
            | 527e3807-4dfe-46dc-a4b8-c2eeba6aab62 |
            | c59ddeae-ba53-4290-8304-61bb59dd2f15 |
            | ebc831fc-ef0b-486a-b8e1-291aecc601ba |
            | 6c37bb86-be69-4a2d-8e7a-dbb6f48f2e29 |
            | 532a2efb-958e-4368-be85-480fe049209d |
            | 3682f081-01be-4c17-a988-28ce54b0999c |
            | 402d6465-f520-4338-a007-ddbcbdb5986c |
        When I request a list of payments
        Then I get an "OK" response
        And the response contains a list of payments with the following IDs:
            | 199aa4d7-1af9-4933-9948-44b98d8ba03e |
            | 233f2a4b-631e-4f0c-a1dc-c33a3792e9b1 |
            | 23ffa122-44aa-4106-b036-877e61848d29 |
            | 29bc7b20-f134-4aee-954d-36fb1daa8c13 |
            | 34d746ce-1f7d-4ea6-958d-28dda4b2c3cb |
            | 3682f081-01be-4c17-a988-28ce54b0999c |
            | 402d6465-f520-4338-a007-ddbcbdb5986c |
            | 527e3807-4dfe-46dc-a4b8-c2eeba6aab62 |
            | 532a2efb-958e-4368-be85-480fe049209d |
            | 5fd8ed2b-407a-4089-ae25-b2fd70a83250 |

    Scenario: List payments with pagination
        Given there are payments with IDs:
            | 29bc7b20-f134-4aee-954d-36fb1daa8c13 |
            | ac6bd520-c5b9-4573-969a-99d46690aa26 |
            | 233f2a4b-631e-4f0c-a1dc-c33a3792e9b1 |
            | f2fdea8a-2358-4459-b126-eda71065d5c5 |
            | d73e8403-c6a3-49eb-b24b-8c441496a5fd |
            | 9c27337b-2c8a-40b1-abab-a51ec8d6561c |
            | 34d746ce-1f7d-4ea6-958d-28dda4b2c3cb |
            | 199aa4d7-1af9-4933-9948-44b98d8ba03e |
            | eee4cb9a-f9c1-4d51-95f3-0647fb174124 |
            | 80f43626-e895-447e-99b6-6b260627f3ca |
            | 23ffa122-44aa-4106-b036-877e61848d29 |
            | 8d05996d-aad6-4ff2-813a-91130a4188f7 |
            | 5fd8ed2b-407a-4089-ae25-b2fd70a83250 |
            | 527e3807-4dfe-46dc-a4b8-c2eeba6aab62 |
            | c59ddeae-ba53-4290-8304-61bb59dd2f15 |
            | ebc831fc-ef0b-486a-b8e1-291aecc601ba |
            | 6c37bb86-be69-4a2d-8e7a-dbb6f48f2e29 |
            | 532a2efb-958e-4368-be85-480fe049209d |
            | 3682f081-01be-4c17-a988-28ce54b0999c |
            | 402d6465-f520-4338-a007-ddbcbdb5986c |
        When I request a list of payments, page 2 with 5 payments per page
        Then I get an "OK" response
        And the response contains a list of payments with the following IDs:
            | 6c37bb86-be69-4a2d-8e7a-dbb6f48f2e29 |
            | 80f43626-e895-447e-99b6-6b260627f3ca |
            | 8d05996d-aad6-4ff2-813a-91130a4188f7 |
            | 9c27337b-2c8a-40b1-abab-a51ec8d6561c |
            | ac6bd520-c5b9-4573-969a-99d46690aa26 |
